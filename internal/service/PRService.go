package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/internal/entity"
	"PR-appointer/internal/repository"
)

type PRService struct {
	prRepo   *repository.PRRepository
	userRepo *repository.UserRepository
	teamRepo *repository.TeamRepository
}

func NewPRService(db *pgxpool.Pool) *PRService {
	return &PRService{
		prRepo:   repository.NewPRRepository(db),
		userRepo: repository.NewUserRepository(db),
		teamRepo: repository.NewTeamRepository(db),
	}
}

func (s *PRService) CreatePR(ctx context.Context, req *entity.PRCreateRequest) (*entity.PRDetailResponse, error) {
	author, err := s.userRepo.GetByID(ctx, req.AuthorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// Получаем команды автора
	teamIDs, err := s.userRepo.GetTeamsByUserID(ctx, req.AuthorID)
	if err != nil {
		return nil, err
	}

	if len(teamIDs) == 0 {
		return nil, errors.New("author is not in any team")
	}

	pr, err := s.prRepo.Create(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("PR already exists")
	}

	// Назначаем ревьюверов из первой команды автора
	teamID := teamIDs[0]
	reviewers, err := s.assignReviewers(ctx, pr.ID, teamID, req.AuthorID)
	if err != nil {
		reviewers = []entity.UserResponse{}
	}

	// Формируем ответ
	var reviewerResponses []entity.UserResponse
	for _, reviewer := range reviewers {
		reviewerResponses = append(reviewerResponses, entity.UserResponse{
			UserID:   reviewer.UserID,
			Username: reviewer.Username,
			IsActive: reviewer.IsActive,
			TeamName: reviewer.TeamName,
		})
	}

	return &entity.PRDetailResponse{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Title,
		Author: entity.UserResponse{
			UserID:   author.UserID,
			Username: author.Username,
			IsActive: author.IsActive,
			TeamName: author.TeamName,
		},
		Status:    string(pr.Status),
		Reviewers: reviewerResponses,
	}, nil
}

func (s *PRService) assignReviewers(ctx context.Context, prID, teamID, authorID int) ([]entity.UserResponse, error) {
	candidates, err := s.teamRepo.GetActiveMembers(ctx, teamID, &authorID)
	if err != nil {
		return nil, err
	}

	// Выбираем до 2 случайных ревьюверов
	numReviewers := min(2, len(candidates))

	// Перемешиваем кандидатов
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	var assignedReviewers []entity.UserResponse
	for i := 0; i < numReviewers; i++ {
		reviewer := candidates[i]
		if err := s.prRepo.AddReviewer(ctx, prID, reviewer.UserID); err != nil {
			continue
		}
		assignedReviewers = append(assignedReviewers, reviewer)
	}

	return assignedReviewers, nil
}

func (s *PRService) MergePR(ctx context.Context, prID int) (*entity.MergedPRResponse, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		slog.Error("Error getting PR", strconv.Itoa(prID), err)
		return nil, err
	}

	// если уже MERGED, просто возвращаем
	if pr.Status == "MERGED" {
		return s.getMergedPRDetails(ctx, pr)
	}

	// Обновляем статус
	pr, err = s.prRepo.UpdateStatus(ctx, prID, "MERGED")
	if err != nil {
		slog.Error("Error updating PR", strconv.Itoa(prID), err)
		return nil, err
	}

	return s.getMergedPRDetails(ctx, pr)
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID int, oldReviewerID int) (*entity.PRDetailResponse, string, error) {
	pr, err := s.validatePRForReassignment(ctx, prID, oldReviewerID)
	if err != nil {
		slog.Error("Error validating PR", strconv.Itoa(prID), err)
		return nil, "", err
	}

	// Получаем старого ревьювера
	_, err = s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		slog.Error("Error user not found", err.Error(), nil)
		return nil, "", errors.New("user not found")
	}

	// Получаем команду старого ревьювера
	teamIDs, err := s.userRepo.GetTeamsByUserID(ctx, oldReviewerID)
	if err != nil || len(teamIDs) == 0 {
		slog.Error("Error getting teamIDs or there is no team", err.Error(), nil)
		return nil, "", errors.New("no active replacement candidate in team")
	}

	// Получаем текущих ревьюверов
	currentReviewers, err := s.prRepo.GetReviewers(ctx, pr.ID)
	if err != nil {
		slog.Error("Error getting current Reviewers ", strconv.Itoa(pr.ID), err.Error())
		return nil, "", err
	}

	// Создаем список ID текущих ревьюверов для исключения
	excludeIDs := make(map[int]bool)
	excludeIDs[pr.AuthorID] = true // Исключаем автора
	for _, r := range currentReviewers {
		excludeIDs[r.UserID] = true
	}

	// Ищем замену из первой команды
	teamID := teamIDs[0]
	candidates, err := s.teamRepo.GetActiveMembers(ctx, teamID, &pr.AuthorID)
	if err != nil {
		slog.Error("Error getting new reviewers", err.Error(), nil)
		return nil, "", err
	}

	// Фильтруем кандидатов
	var availableCandidates []entity.UserResponse
	for _, candidate := range candidates {
		if !excludeIDs[candidate.UserID] {
			availableCandidates = append(availableCandidates, candidate)
		}
	}

	if len(availableCandidates) == 0 {
		slog.Warn("No candidates")
		return nil, "", errors.New("no active replacement candidate in team")
	}

	// Выбираем случайного кандидата
	newReviewer := availableCandidates[rand.Intn(len(availableCandidates))]

	// Удаляем старого ревьювера
	if err = s.prRepo.RemoveReviewer(ctx, pr.ID, oldReviewerID); err != nil {
		slog.Error("Error removing old reviewer", strconv.Itoa(pr.ID), err.Error())
		return nil, "", err
	}

	// Добавляем нового ревьювера
	if err = s.prRepo.AddReviewer(ctx, pr.ID, newReviewer.UserID); err != nil {
		slog.Error("Error adding newReviewer", strconv.Itoa(pr.ID), err.Error())
		return nil, "", err
	}

	// Формируем ответ
	prDetails, err := s.getPRDetails(ctx, pr)
	if err != nil {
		slog.Error("Error getting prDetails", strconv.Itoa(pr.ID), err.Error())
		return nil, "", err
	}

	return prDetails, fmt.Sprintf("%d", newReviewer.UserID), nil
}

func (s *PRService) validatePRForReassignment(ctx context.Context, prID, reviewerID int) (*entity.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == "MERGED" {
		return nil, errors.New("cannot reassign on merged PR")
	}

	isAssigned, err := s.prRepo.IsReviewerAssigned(ctx, prID, reviewerID)
	if err != nil {
		return nil, err
	}

	if !isAssigned {
		return nil, errors.New("reviewer is not assigned to this PR")
	}

	return pr, nil
}

func (s *PRService) getPRDetails(ctx context.Context, pr *entity.PullRequest) (*entity.PRDetailResponse, error) {
	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.prRepo.GetReviewers(ctx, pr.ID)
	if err != nil {
		reviewers = []entity.UserResponse{}
	}

	var reviewerResponses []entity.UserResponse
	for _, reviewer := range reviewers {
		reviewerResponses = append(reviewerResponses, entity.UserResponse{
			UserID:   reviewer.UserID,
			Username: reviewer.Username,
			IsActive: reviewer.IsActive,
			TeamName: reviewer.TeamName,
		})
	}

	return &entity.PRDetailResponse{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Title,
		Author: entity.UserResponse{
			UserID:   author.UserID,
			Username: author.Username,
			IsActive: author.IsActive,
			TeamName: author.TeamName,
		},
		Status:    string(pr.Status),
		Reviewers: reviewerResponses,
	}, nil
}

func (s *PRService) getMergedPRDetails(ctx context.Context, pr *entity.PullRequest) (*entity.MergedPRResponse, error) {
	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.prRepo.GetReviewers(ctx, pr.ID)
	if err != nil {
		reviewers = []entity.UserResponse{}
	}

	var reviewerResponses []entity.UserResponse
	for _, reviewer := range reviewers {
		reviewerResponses = append(reviewerResponses, entity.UserResponse{
			UserID:   reviewer.UserID,
			Username: reviewer.Username,
			IsActive: reviewer.IsActive,
			TeamName: reviewer.TeamName,
		})
	}

	return &entity.MergedPRResponse{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Title,
		Author: entity.UserResponse{
			UserID:   author.UserID,
			Username: author.Username,
			IsActive: author.IsActive,
			TeamName: author.TeamName,
		},
		Status:    string(pr.Status),
		Reviewers: reviewerResponses,
		MergedAt:  pr.MergedAt,
	}, nil
}
