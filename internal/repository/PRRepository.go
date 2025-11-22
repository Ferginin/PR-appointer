package repository

import (
	"PR-appointer/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRRepository struct {
	db *pgxpool.Pool
}

func NewPRRepository(db *pgxpool.Pool) *PRRepository {
	return &PRRepository{db: db}
}

func (r *PRRepository) Create(ctx context.Context, id int, title string, authorID int) (*entity.PullRequest, error) {
	query := `
		INSERT INTO pull_requests (id, title, author_id, status)
		VALUES ($1, $2, $3, 'OPEN')
		RETURNING id, title, author_id, status, created_at, updated_at
	`

	pr := entity.PullRequest{}
	err := r.db.QueryRow(ctx, query, id, title, authorID).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return &pr, nil
}

func (r *PRRepository) GetByID(ctx context.Context, prID int) (*entity.PullRequest, error) {
	query := `
		SELECT id, title, author_id, status, created_at, updated_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	pr := entity.PullRequest{}
	err := r.db.QueryRow(ctx, query, prID).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	return &pr, nil
}

func (r *PRRepository) UpdateStatus(ctx context.Context, prID int, status string) (*entity.PullRequest, error) {
	query := `
		UPDATE pull_requests
		SET status = $1, updated_at = CURRENT_TIMESTAMP, merged_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, title, author_id, status, created_at, updated_at, merged_at
	`

	pr := entity.PullRequest{}
	err := r.db.QueryRow(ctx, query, status, prID).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("PR not found")
		}
		return nil, fmt.Errorf("failed to update PR status: %w", err)
	}

	return &pr, nil
}

func (r *PRRepository) AddReviewer(ctx context.Context, prID int, reviewerID int) error {
	query := `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to add reviewer: %w", err)
	}

	return nil
}

func (r *PRRepository) RemoveReviewer(ctx context.Context, prID int, reviewerID int) error {
	query := `
		DELETE FROM pr_reviewers
		WHERE pr_id = $1 AND reviewer_id = $2
	`

	_, err := r.db.Exec(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}

	return nil
}

func (r *PRRepository) GetReviewers(ctx context.Context, prID int) ([]entity.UserResponse, error) {
	query := `
		SELECT u.id, u.username, u.is_active, t.name
		FROM users u
		JOIN pr_reviewers pr ON u.id = pr.reviewer_id
		JOIN team_members on u.id = team_members.user_id
		JOIN teams t on team_members.team_id = t.id
		WHERE pr.pr_id = $1
		ORDER BY pr.assigned_at
	`

	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviewers: %w", err)
	}

	return ScanUserResponses(ctx, rows) // Используем общую функцию
}

func (r *PRRepository) IsReviewerAssigned(ctx context.Context, prID int, reviewerID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers 
			WHERE pr_id = $1 AND reviewer_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, prID, reviewerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check reviewer assignment: %w", err)
	}

	return exists, nil
}

func (r *PRRepository) GetPRsByReviewer(ctx context.Context, reviewerID int) ([]entity.PullRequest, error) {
	query := `
		SELECT DISTINCT pr.id, pr.title, pr.author_id, pr.status, pr.created_at, pr.updated_at
		FROM pull_requests pr
		INNER JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []entity.PullRequest
	for rows.Next() {
		pr := entity.PullRequest{}
		err := rows.Scan(
			&pr.ID,
			&pr.Title,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
