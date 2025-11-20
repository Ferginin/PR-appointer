package repository

import (
	"PR-appointer/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username string, isActive bool) (*entity.User, error) {
	query := `
		INSERT INTO users (username, is_active)
		VALUES ($1, $2)
		RETURNING id, username, is_active, created_at, updated_at
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, username, isActive).Scan(
		&user.ID,
		&user.Username,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID int) (*entity.UserResponse, error) {
	query := `
		SELECT users.id, username, is_active, teams.name FROM users
		JOIN team_members on team_members.user_id = users.id
		JOIN teams on teams.id = team_members.team_id
		WHERE users.id = $1
	`

	var user entity.UserResponse
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.IsActive,
		&user.TeamName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, username, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // User doesn't exist
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID int, isActive bool) (*entity.UserResponse, error) {
	query := `
		UPDATE users
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	r.db.QueryRow(ctx, query, isActive, userID)

	query = `
		SELECT users.id, username, teams.name, is_active FROM users
		JOIN team_members on team_members.user_id = users.id
		JOIN teams on teams.id = team_members.team_id
		WHERE users.id = $1
	`

	var user entity.UserResponse
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepository) Upsert(ctx context.Context, username string, isActive bool) (*entity.User, error) {
	query := `
		INSERT INTO users (username, is_active)
		VALUES ($1, $2)
		ON CONFLICT (username) 
		DO UPDATE SET is_active = EXCLUDED.is_active, updated_at = CURRENT_TIMESTAMP
		RETURNING id, username, is_active, created_at, updated_at
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, username, isActive).Scan(
		&user.ID,
		&user.Username,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetTeamsByUserID(ctx context.Context, userID int) ([]int, error) {
	query := `
		SELECT team_id
		FROM team_members
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user teams: %w", err)
	}
	defer rows.Close()

	var teamIDs []int
	for rows.Next() {
		var teamID int
		if err := rows.Scan(&teamID); err != nil {
			return nil, fmt.Errorf("failed to scan team id: %w", err)
		}
		teamIDs = append(teamIDs, teamID)
	}

	return teamIDs, nil
}
