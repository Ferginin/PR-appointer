package service

import (
	"PR-appointer/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"math/rand"

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

	// Проверяем существование автора
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

	// Создаем PR
	pr, err := s.prRepo.Create(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("PR already exists")
	}

	// Назначаем ревьюверов из первой команды автора
	teamID := teamIDs[0]
	reviewers, err := s.assignReviewers(ctx, pr.ID, teamID, req.AuthorID)
	if err != nil {
		// PR создан, но ревьюверов назначить не удалось - это нормально
		reviewers = []entity.User{}
	}

	// Формируем ответ
	var reviewerResponses []entity.UserResponse
	for _, reviewer := range reviewers {
		reviewerResponses = append(reviewerResponses, entity.UserResponse{
			UserID:   reviewer.ID,
			Username: reviewer.Username,
			IsActive: reviewer.IsActive,
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

func (s *PRService) assignReviewers(ctx context.Context, prID, teamID, authorID int) ([]entity.User, error) {
	// Получаем активных участников команды, исключая автора
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

	var assignedReviewers []entity.User
	for i := 0; i < numReviewers; i++ {
		reviewer := candidates[i]
		if err := s.prRepo.AddReviewer(ctx, prID, reviewer.ID); err != nil {
			// Игнорируем ошибку и продолжаем
			continue
		}
		assignedReviewers = append(assignedReviewers, reviewer)
	}

	return assignedReviewers, nil
}

func (s *PRService) MergePR(ctx context.Context, prID int) (*entity.PRDetailResponse, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	// Идемпотентная операция - если уже MERGED, просто возвращаем
	if pr.Status == "MERGED" {
		return s.getPRDetails(ctx, pr)
	}

	// Обновляем статус
	pr, err = s.prRepo.UpdateStatus(ctx, prID, "MERGED")
	if err != nil {
		return nil, err
	}

	return s.getPRDetails(ctx, pr)
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID int, oldReviewerID int) (*entity.PRDetailResponse, string, error) {
	// Получаем PR
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		slog.Error("Error getting PR Details", slog.String("error", err.Error()))
		return nil, "", err
	}

	// Проверяем статус
	if pr.Status == "MERGED" {
		slog.Warn("PR is merged, skip reassign reviewer")
		return nil, "", errors.New("cannot reassign on merged PR")
	}

	// Проверяем, назначен ли ревьювер
	isAssigned, err := s.prRepo.IsReviewerAssigned(ctx, pr.ID, oldReviewerID)
	if err != nil {
		slog.Error("Error checking if PR is assigned", slog.String("error", err.Error()))
		return nil, "", err
	}

	if !isAssigned {
		slog.Warn("PR is not assigned on reviewer")
		return nil, "", errors.New("reviewer is not assigned to this PR")
	}

	// Получаем старого ревьювера
	_, err = s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		slog.Error("Error user not found", err.Error())
		return nil, "", errors.New("user not found")
	}

	// Получаем команды старого ревьювера
	teamIDs, err := s.userRepo.GetTeamsByUserID(ctx, oldReviewerID)
	if err != nil || len(teamIDs) == 0 {
		slog.Error("Error getting teamIDs or there is no team", err)
		return nil, "", errors.New("no active replacement candidate in team")
	}

	// Получаем текущих ревьюверов
	currentReviewers, err := s.prRepo.GetReviewers(ctx, pr.ID)
	if err != nil {
		slog.Error("Error getting current Reviewers ", pr.ID, err.Error())
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
		slog.Error("Error getting new reviewers", err.Error())
		return nil, "", err
	}

	// Фильтруем кандидатов
	var availableCandidates []entity.User
	for _, candidate := range candidates {
		if !excludeIDs[candidate.ID] {
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
	if err := s.prRepo.RemoveReviewer(ctx, pr.ID, oldReviewerID); err != nil {
		slog.Error("Error removing old reviewer", pr.ID, err.Error())
		return nil, "", err
	}

	// Добавляем нового ревьювера
	if err := s.prRepo.AddReviewer(ctx, pr.ID, newReviewer.ID); err != nil {
		slog.Error("Error adding newReviewer", pr.ID, err.Error())
		return nil, "", err
	}

	// Формируем ответ
	prDetails, err := s.getPRDetails(ctx, pr)
	if err != nil {
		slog.Error("Error getting prDetails", pr, err.Error())
		return nil, "", err
	}

	return prDetails, fmt.Sprintf("%d", newReviewer.ID), nil
}

func (s *PRService) getPRDetails(ctx context.Context, pr *entity.PullRequest) (*entity.PRDetailResponse, error) {
	// Получаем автора
	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	// Получаем ревьюверов
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
