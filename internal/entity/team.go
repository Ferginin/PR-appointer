package entity

import "time"

type Team struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TeamCreateRequest struct {
	TeamName string         `json:"team_name"`
	Members  []UserResponse `json:"members"`
}

type TeamResponse struct {
	TeamName string         `json:"team_name"`
	Members  []UserResponse `json:"members"`
}
