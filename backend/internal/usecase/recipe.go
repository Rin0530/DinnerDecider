package usecase

import (
	"context"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
)

// RecipeUsecase defines the business logic interface for recipe operations
type RecipeUsecase interface {
	// GetRecipeSuggestion generates recipe suggestions based on available ingredients
	GetRecipeSuggestion(ctx context.Context) (*domain.RecipeResponse, error)
}
