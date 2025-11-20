package service

import (
	"PR-appointer/internal/entity"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"

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
	// Проверяем, существует ли команда
	existingTeam, _ := s.teamRepo.GetByName(ctx, req.TeamName)
	if existingTeam != nil {
		return nil, errors.New("team already exists")
	}

	// Создаем команду
	team, err := s.teamRepo.Create(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}

	// Создаем/обновляем пользователей и добавляем в команду
	var members []entity.UserResponse
	for _, member := range req.Members {
		// Upsert пользователя
		user, err := s.userRepo.Upsert(ctx, member.Username, member.IsActive)
		if err != nil {
			return nil, err
		}

		// Добавляем в команду
		if err := s.teamRepo.AddMember(ctx, team.ID, user.ID); err != nil {
			return nil, err
		}

		members = append(members, entity.UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
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

	// Получаем участников команды
	members, err := s.teamRepo.GetMembers(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	var memberResponses []entity.UserResponse
	for _, member := range members {
		memberResponses = append(memberResponses, entity.UserResponse{
			UserID:   member.ID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return &entity.TeamResponse{
		TeamName: team.Name,
		Members:  memberResponses,
	}, nil
}
