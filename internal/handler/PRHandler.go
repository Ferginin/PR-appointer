package handler

import (
	"PR-appointer/config"
	"PR-appointer/internal/entity"
	"PR-appointer/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(cfg *config.Config, db *pgxpool.Pool) *PRHandler {
	return &PRHandler{
		prService: service.NewPRService(db),
	}
}

// CreatePR godoc
// @Summary Create PR with auto-assigned reviewers
// @Description Create PR and automatically assign up to 2 reviewers from author's team
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body entity.PRCreateRequest true "PR data"
// @Success 201 {object} entity.PRDetailResponse
// @Failure 404 {object} APIError
// @Failure 409 {object} APIError
// @Router /pullRequest/create [post]
func (h *PRHandler) CreatePR(c *gin.Context) {
	var req entity.PRCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	pr, err := h.prService.CreatePR(c.Request.Context(), &req)
	if err != nil {
		// Check specific error types
		if err.Error() == "PR already exists" {
			c.JSON(http.StatusConflict, newAPIError(ErrCodePRExists, "PR id already exists"))
			return
		}
		if err.Error() == "author not found" || err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, newAPIError(ErrCodeNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

// MergePR godoc
// @Summary Mark PR as merged
// @Description Set PR status to MERGED (idempotent operation)
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body entity.UpdatePRStatusRequest true "PR ID"
// @Success 200 {object} entity.PRDetailResponse
// @Failure 404 {object} APIError
// @Router /pullRequest/merge [post]
func (h *PRHandler) MergePR(c *gin.Context) {
	var req entity.UpdatePRStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	pr, err := h.prService.MergePR(c.Request.Context(), req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, newAPIError(ErrCodeNotFound, "PR not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

// ReassignReviewer godoc
// @Summary Reassign reviewer
// @Description Replace one reviewer with another from same team
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body entity.ReassignReviewerRequest true "Reassignment data"
// @Success 200 {object} entity.ReassignHandlerResponse
// @Failure 404 {object} APIError
// @Failure 409 {object} APIError
// @Router /pullRequest/reassign [post]
func (h *PRHandler) ReassignReviewer(c *gin.Context) {
	var req entity.ReassignReviewerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		// Check specific error types
		switch err.Error() {
		case "PR not found", "user not found":
			c.JSON(http.StatusNotFound, newAPIError(ErrCodeNotFound, err.Error()))
			return
		case "cannot reassign on merged PR":
			c.JSON(http.StatusConflict, newAPIError(ErrCodePRMerged, err.Error()))
			return
		case "reviewer is not assigned to this PR":
			c.JSON(http.StatusConflict, newAPIError(ErrCodeNotAssigned, err.Error()))
			return
		case "no active replacement candidate in team":
			c.JSON(http.StatusConflict, newAPIError(ErrCodeNoCandidate, err.Error()))
			return
		default:
			c.JSON(http.StatusBadRequest, newAPIError(ErrCodeNotFound, err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}
