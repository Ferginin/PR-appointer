package handler

import (
	"PR-appointer/config"
	"PR-appointer/internal/entity"
	"PR-appointer/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(cfg *config.Config, db *pgxpool.Pool) *TeamHandler {
	return &TeamHandler{
		teamService: service.NewTeamService(db),
	}
}

// AddTeam godoc
// @Summary Create team with members
// @Description Create team and create/update users in it
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body entity.TeamCreateRequest true "Team data"
// @Success 201 {object} entity.TeamResponse
// @Failure 400 {object} APIError
// @Router /team/add [post]
func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req entity.TeamCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	team, err := h.teamService.CreateTeam(c.Request.Context(), &req)
	if err != nil {
		// Check if team already exists
		if err.Error() == "team already exists" {
			c.JSON(http.StatusBadRequest, newAPIError(ErrCodeTeamExists, "team_name already exists"))
			return
		}
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": team})
}

// GetTeam godoc
// @Summary Get team with members
// @Description Get team by name
// @Tags Teams
// @Produce json
// @Param team_name query string true "Team name"
// @Success 200 {object} entity.TeamResponse
// @Failure 404 {object} APIError
// @Router /team/get [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, "team_name is required"))
		return
	}

	team, err := h.teamService.GetTeamByName(c.Request.Context(), teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, newAPIError(ErrCodeNotFound, "team not found"))
		return
	}

	c.JSON(http.StatusOK, team)
}
