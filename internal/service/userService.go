package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/internal/entity"
	"PR-appointer/internal/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
	prRepo   *repository.PRRepository
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		UserRepo: repository.NewUserRepository(db),
		prRepo:   repository.NewPRRepository(db),
	}
}

func (s *UserService) SetStatus(ctx context.Context, userID int, isActive bool) (*entity.UserResponse, error) {
	user, err := s.UserRepo.UpdateStatus(ctx, userID, isActive)
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
	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

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
