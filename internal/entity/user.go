package entity

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username" binding:"required"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRequest struct {
	UserID   int  `json:"user_id" binding:"required"`
	IsActive bool `json:"is_active,omitempty"`
}

type UserResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name" db:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserReviewsResponse struct {
	UserID       int         `json:"user_id"`
	Username     string      `json:"username"`
	PullRequests []PRSummary `json:"pull_requests"`
}

type PRSummary struct {
	PullRequestID   int    `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        int    `json:"author_id"`
	Status          string `json:"status"`
}
