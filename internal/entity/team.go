package entity

import "time"

type Team struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TeamMember представляет связь пользователя с командой
type TeamMember struct {
	ID        int       `json:"id" db:"id"`
	TeamID    int       `json:"team_id" db:"team_id" binding:"required"`
	UserID    int       `json:"user_id" db:"user_id" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateTeamRequest - запрос на создание команды
type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// AddTeamMemberRequest - запрос на добавление пользователя в команду
type AddTeamMemberRequest struct {
	UserID int `json:"user_id" binding:"required,min=1"`
}

type TeamCreateRequest struct {
	TeamName string         `json:"team_name"`
	Members  []UserResponse `json:"members"`
}

type TeamResponse struct {
	TeamName string         `json:"team_name"`
	Members  []UserResponse `json:"members"`
}
