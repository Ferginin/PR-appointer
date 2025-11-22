package entity

import (
	"database/sql"
	"time"
)

type PullRequest struct {
	ID        int          `json:"id" db:"id"`
	Title     string       `json:"title" db:"title" binding:"required"`
	AuthorID  int          `json:"author_id" db:"author_id" binding:"required"`
	Status    string       `json:"status" db:"status"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	MergedAt  sql.NullTime `json:"merged_at" db:"merged_at"`
}

type UpdatePRStatusRequest struct {
	PullRequestID int `json:"pull_request_id" binding:"required"`
}

type ReassignReviewerRequest struct {
	PullRequestID int `json:"pull_request_id" binding:"required"`
	OldReviewerID int `json:"old_reviewer_id" binding:"required,min=1"`
}

type PRCreateRequest struct {
	PullRequestID   int    `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        int    `json:"author_id"`
}

type PRDetailResponse struct {
	PullRequestID   int            `json:"pull_request_id"`
	PullRequestName string         `json:"pull_request_name"`
	Author          UserResponse   `json:"author"`
	Status          string         `json:"status"`
	Reviewers       []UserResponse `json:"reviewers"`
}

type MergedPRResponse struct {
	PullRequestID   int            `json:"pull_request_id"`
	PullRequestName string         `json:"pull_request_name"`
	Author          UserResponse   `json:"author"`
	Status          string         `json:"status"`
	Reviewers       []UserResponse `json:"reviewers"`
	MergedAt        sql.NullTime   `json:"merged_at" db:"merged_at"`
}

type ReassignHandlerResponse struct {
	Pr         PRDetailResponse `json:"pr"`
	ReplacedBy string           `json:"replaced_by"`
}
