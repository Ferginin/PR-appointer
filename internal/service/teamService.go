package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/internal/entity"
	"PR-appointer/internal/repository"
)

type TeamService struct {
	teamRepo *repository.TeamRepository
	userRepo *repository.UserRepository
}

func NewTeamService(db *pgxpool.Pool) *TeamService {
	return &TeamService{
		teamRepo: repository.NewTeamRepository(db),
		userRepo: repository.NewUserRepository(db),
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, req *entity.TeamCreateRequest) (*entity.TeamResponse, error) {
	existingTeam, err := s.teamRepo.GetByName(ctx, req.TeamName)
	if existingTeam != nil || err != nil {
		return nil, errors.New("team already exists")
	}

	team, err := s.teamRepo.Create(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}

	// Создаем/обновляем пользователей и добавляем в команду
	var members []entity.UserResponse
	for _, member := range req.Members {
		user, err := s.userRepo.Upsert(ctx, member.Username, member.IsActive)
		if err != nil {
			return nil, err
		}

		if err := s.teamRepo.AddMember(ctx, team.ID, user.UserID); err != nil {
			return nil, err
		}

		members = append(members, entity.UserResponse{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
			TeamName: user.TeamName,
		})
	}

	return &entity.TeamResponse{
		TeamName: team.Name,
		Members:  members,
	}, nil
}

func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*entity.TeamResponse, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	members, err := s.teamRepo.GetMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	var memberResponses []entity.UserResponse
	for _, member := range members {
		memberResponses = append(memberResponses, entity.UserResponse{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: member.TeamName,
		})
	}

	return &entity.TeamResponse{
		TeamName: team.Name,
		Members:  memberResponses,
	}, nil
}
