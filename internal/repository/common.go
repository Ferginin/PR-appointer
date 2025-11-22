package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"PR-appointer/internal/entity"
)

func ScanUserResponses(ctx context.Context, rows pgx.Rows) ([]entity.UserResponse, error) {
	defer rows.Close()

	var users []entity.UserResponse
	for rows.Next() {
		user := entity.UserResponse{}
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.IsActive,
			&user.TeamName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
