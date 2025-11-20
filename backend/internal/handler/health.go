package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db            *sqlx.DB
	ollamaService service.OllamaService
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler(db *sqlx.DB, ollamaService service.OllamaService) *HealthHandler {
	return &HealthHandler{
		db:            db,
		ollamaService: ollamaService,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health handles GET /health - basic application health check
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

// HealthDB handles GET /health/db - database connection health check
func (h *HealthHandler) HealthDB(c *gin.Context) {
	// Create a context with timeout for the ping operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Ping the database
	if err := h.db.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:  "error",
			Message: fmt.Sprintf("Database connection failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

// HealthOllama handles GET /health/ollama - Ollama API health check
func (h *HealthHandler) HealthOllama(c *gin.Context) {
	// Create a context with timeout for the health check
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Try to generate a simple recipe suggestion with empty ingredients
	// This will test if Ollama API is reachable and responding
	_, err := h.ollamaService.GenerateRecipeSuggestion(ctx, nil)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:  "error",
			Message: fmt.Sprintf("Ollama API connection failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}
