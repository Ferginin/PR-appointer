package service

import (
	"PR-appointer/internal/entity"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
	prRepo   *repository.PRRepository
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(db),
		prRepo:   repository.NewPRRepository(db),
	}
}

func (s *UserService) SetStatus(ctx context.Context, userID int, isActive bool) (*entity.UserResponse, error) {
	user, err := s.userRepo.UpdateStatus(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return &entity.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}, nil
}

func (s *UserService) GetUserReviews(ctx context.Context, userID int) (*entity.UserReviewsResponse, error) {
	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем PR, назначенные на пользователя
	prs, err := s.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	var prSummaries []entity.PRSummary
	for _, pr := range prs {
		prSummaries = append(prSummaries, entity.PRSummary{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Title,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	return &entity.UserReviewsResponse{
		UserID:       user.UserID,
		Username:     user.Username,
		PullRequests: prSummaries,
	}, nil
}
