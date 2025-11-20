package usecase

import (
	"context"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
)

// IngredientUsecase defines the business logic interface for ingredient operations
type IngredientUsecase interface {
	// CreateIngredient creates a new ingredient
	CreateIngredient(ctx context.Context, req CreateIngredientRequest) (*domain.Ingredient, error)

	// GetAllIngredients retrieves all ingredients
	GetAllIngredients(ctx context.Context) ([]*domain.Ingredient, error)

	// GetIngredientByID retrieves an ingredient by ID
	GetIngredientByID(ctx context.Context, id int64) (*domain.Ingredient, error)

	// UpdateIngredient updates an existing ingredient
	UpdateIngredient(ctx context.Context, id int64, req UpdateIngredientRequest) (*domain.Ingredient, error)

	// DeleteIngredient deletes an ingredient by ID
	DeleteIngredient(ctx context.Context, id int64) error
}
