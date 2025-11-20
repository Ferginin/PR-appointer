package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct{}

// HealthCheck godoc
// @Summary Health check
// @Description Check service health status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
