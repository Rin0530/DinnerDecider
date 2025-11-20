package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// respondWithError sends a standardized error response
func respondWithError(c *gin.Context, statusCode int, errorType string, message string) {
	c.JSON(statusCode, usecase.ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

// handleError maps common errors to appropriate HTTP status codes and responses
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Handle sql.ErrNoRows (resource not found)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(c, http.StatusNotFound, "not_found", "Resource not found")
		return
	}

	// Default to internal server error
	respondWithError(c, http.StatusInternalServerError, "internal_error", err.Error())
}

// respondBadRequest sends a 400 Bad Request response
func respondBadRequest(c *gin.Context, message string) {
	respondWithError(c, http.StatusBadRequest, "validation_error", message)
}

// respondNotFound sends a 404 Not Found response
func respondNotFound(c *gin.Context, message string) {
	respondWithError(c, http.StatusNotFound, "not_found", message)
}

// respondInternalError sends a 500 Internal Server Error response
func respondInternalError(c *gin.Context, message string) {
	respondWithError(c, http.StatusInternalServerError, "internal_error", message)
}

// respondServiceUnavailable sends a 503 Service Unavailable response
func respondServiceUnavailable(c *gin.Context, message string) {
	respondWithError(c, http.StatusServiceUnavailable, "service_unavailable", message)
}
