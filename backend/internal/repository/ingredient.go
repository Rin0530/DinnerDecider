package repository

import (
	"context"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
)

// IngredientRepository defines the interface for ingredient data access
type IngredientRepository interface {
	// Create inserts a new ingredient into the database
	Create(ctx context.Context, ingredient *domain.Ingredient) error

	// GetAll retrieves all ingredients from the database
	GetAll(ctx context.Context) ([]*domain.Ingredient, error)

	// GetByID retrieves a single ingredient by its ID
	GetByID(ctx context.Context, id int64) (*domain.Ingredient, error)

	// Update modifies an existing ingredient in the database
	Update(ctx context.Context, ingredient *domain.Ingredient) error

	// Delete removes an ingredient from the database by its ID
	Delete(ctx context.Context, id int64) error
}
