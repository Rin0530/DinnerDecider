package service

import (
	"context"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
)

// OllamaService defines the interface for interacting with Ollama API
type OllamaService interface {
	// GenerateRecipeSuggestion generates recipe suggestions based on available ingredients
	GenerateRecipeSuggestion(ctx context.Context, ingredients []*domain.Ingredient) (*domain.RecipeResponse, error)
}
