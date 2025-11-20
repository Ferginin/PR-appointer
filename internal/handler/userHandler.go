package handler

import (
	"net/http"
	"strconv"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"PR-appointer/config"
	"PR-appointer/internal/entity"
	"PR-appointer/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(cfg *config.Config, db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(db),
	}
}

// SetStatus godoc
// @Summary Set user active status
// @Description Update user's is_active flag
// @Tags Users
// @Accept json
// @Produce json
// @Param request body entity.UserRequest true "User active status"
// @Success 200 {object} entity.UserResponse
// @Failure 404 {object} APIError
// @Router /users/setIsActive [post]
func (h *UserHandler) SetStatus(c *gin.Context) {
	var req entity.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	user, err := h.userService.SetStatus(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, newAPIError(ErrCodeNotFound, "user not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUserReviews godoc
// @Summary Get PRs assigned to user as reviewer
// @Description Get list of PRs where user is assigned as reviewer
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} entity.UserReviewsResponse
// @Router /users/getReview [get]
func (h *UserHandler) GetUserReviews(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, "user_id is required"))
		return
	}

	reviews, err := h.userService.GetUserReviews(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, reviews)
}
