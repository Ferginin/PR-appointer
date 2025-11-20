package entity

import "time"

// PRStatus представляет статус Pull Request
type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

// PullRequest представляет Pull Request
type PullRequest struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title" binding:"required"`
	AuthorID  int       `json:"author_id" db:"author_id" binding:"required"`
	Status    PRStatus  `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

//// PRReviewer представляет назначенного ревьювера на PR
//type PRReviewer struct {
//	ID         int       `json:"id" db:"id"`
//	PRID       int       `json:"pr_id" db:"pr_id" binding:"required"`
//	ReviewerID int       `json:"reviewer_id" db:"reviewer_id" binding:"required"`
//	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
//}
//
//// CreatePRRequest - запрос на создание PR
//type CreatePRRequest struct {
//	Title    string `json:"title" binding:"required,min=1,max=500"`
//	AuthorID int    `json:"author_id" binding:"required,min=1"`
//}
//

// UpdatePRStatusRequest - запрос на обновление статуса PR
type UpdatePRStatusRequest struct {
	PullRequestID int `json:"pull_request_id" binding:"required"`
}

// ReassignReviewerRequest - запрос на переназначение ревьювера
type ReassignReviewerRequest struct {
	PullRequestID int `json:"pull_request_id" binding:"required"`
	OldReviewerID int `json:"old_reviewer_id" binding:"required,min=1"`
}

type ReassignResponse struct {
	PRDetailResponse
	replacedBy UserResponse
}

//// PRResponse - ответ с информацией о PR
//type PRResponse struct {
//	ID        int       `json:"id"`
//	Title     string    `json:"title"`
//	Author    User      `json:"author"`
//	Status    PRStatus  `json:"status"`
//	Reviewers []User    `json:"reviewers"`
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}

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

type ReassignHandlerResponse struct {
	Pr         PRDetailResponse `json:"pr"`
	ReplacedBy string           `json:"replaced_by"`
}
