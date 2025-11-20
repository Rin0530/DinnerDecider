package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/internal/repository"
)

// ingredientUsecase implements the IngredientUsecase interface
type ingredientUsecase struct {
	repo repository.IngredientRepository
}

// NewIngredientUsecase creates a new instance of IngredientUsecase
func NewIngredientUsecase(repo repository.IngredientRepository) IngredientUsecase {
	return &ingredientUsecase{
		repo: repo,
	}
}

// CreateIngredient creates a new ingredient
func (u *ingredientUsecase) CreateIngredient(ctx context.Context, req CreateIngredientRequest) (*domain.Ingredient, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	// Create ingredient domain model
	ingredient := &domain.Ingredient{
		Name:      req.Name,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Parse purchase date if provided
	if req.PurchaseDate != nil && *req.PurchaseDate != "" {
		purchaseDate, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid purchase_date format: %w", err)
		}
		ingredient.PurchaseDate = &purchaseDate
	}

	// Save to repository
	if err := u.repo.Create(ctx, ingredient); err != nil {
		return nil, fmt.Errorf("failed to create ingredient: %w", err)
	}

	return ingredient, nil
}

// GetAllIngredients retrieves all ingredients
func (u *ingredientUsecase) GetAllIngredients(ctx context.Context) ([]*domain.Ingredient, error) {
	ingredients, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ingredients: %w", err)
	}

	// Return empty slice instead of nil if no ingredients found
	if ingredients == nil {
		return []*domain.Ingredient{}, nil
	}

	return ingredients, nil
}

// GetIngredientByID retrieves an ingredient by ID
func (u *ingredientUsecase) GetIngredientByID(ctx context.Context, id int64) (*domain.Ingredient, error) {
	ingredient, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredient by id: %w", err)
	}

	return ingredient, nil
}

// UpdateIngredient updates an existing ingredient
func (u *ingredientUsecase) UpdateIngredient(ctx context.Context, id int64, req UpdateIngredientRequest) (*domain.Ingredient, error) {
	// Get existing ingredient
	ingredient, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredient: %w", err)
	}

	if ingredient == nil {
		return nil, errors.New("ingredient not found")
	}

	// Update fields if provided
	if req.Name != nil {
		ingredient.Name = *req.Name
	}

	if req.Quantity != nil {
		ingredient.Quantity = *req.Quantity
	}

	if req.PurchaseDate != nil {
		if *req.PurchaseDate == "" {
			// Clear purchase date
			ingredient.PurchaseDate = nil
		} else {
			purchaseDate, err := time.Parse("2006-01-02", *req.PurchaseDate)
			if err != nil {
				return nil, fmt.Errorf("invalid purchase_date format: %w", err)
			}
			ingredient.PurchaseDate = &purchaseDate
		}
	}

	// Update timestamp
	ingredient.UpdatedAt = time.Now()

	// Save to repository
	if err := u.repo.Update(ctx, ingredient); err != nil {
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}

	return ingredient, nil
}

// DeleteIngredient deletes an ingredient by ID
func (u *ingredientUsecase) DeleteIngredient(ctx context.Context, id int64) error {
	// Check if ingredient exists
	ingredient, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get ingredient: %w", err)
	}

	if ingredient == nil {
		return errors.New("ingredient not found")
	}

	// Delete from repository
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	return nil
}
