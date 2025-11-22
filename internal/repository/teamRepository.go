package repository

import (
	"PR-appointer/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, name string) (*entity.Team, error) {
	query := `
		INSERT INTO teams (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	team := entity.Team{}
	err := r.db.QueryRow(ctx, query, name).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return &team, nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (*entity.Team, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM teams
		WHERE name = $1
	`

	team := entity.Team{}
	err := r.db.QueryRow(ctx, query, name).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("team not found")
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

func (r *TeamRepository) GetByID(ctx context.Context, teamID int) (*entity.Team, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	team := entity.Team{}
	err := r.db.QueryRow(ctx, query, teamID).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("team not found")
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, teamID, userID int) error {
	query := `
		INSERT INTO team_members (team_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (team_id, user_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetMembers(ctx context.Context, teamID int) ([]entity.UserResponse, error) {
	query := `
		SELECT u.id, u.username, u.is_active, t.name
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		JOIN teams t on tm.team_id = t.id
		WHERE tm.team_id = $1
		ORDER BY u.id
	`

	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query team members: %w", err)
	}

	return ScanUserResponses(ctx, rows) // Используем общую функцию
}

func (r *TeamRepository) GetActiveMembers(ctx context.Context, teamID int, excludeUserID *int) ([]entity.UserResponse, error) {
	query := `
		SELECT u.id, u.username, u.is_active, t.name
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		JOIN teams t on tm.team_id = t.id
		WHERE tm.team_id = $1 AND u.is_active = TRUE
	`

	args := []interface{}{teamID}

	if excludeUserID != nil {
		query += ` AND u.id != $2`
		args = append(args, *excludeUserID)
	}

	query += ` ORDER BY u.id`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query active team members: %w", err)
	}
	defer rows.Close()

	var members []entity.UserResponse
	for rows.Next() {
		user := entity.UserResponse{}
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.IsActive,
			&user.TeamName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}
		members = append(members, user)
	}

	return members, nil
}
