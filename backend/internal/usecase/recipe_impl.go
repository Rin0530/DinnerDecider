package usecase

import (
	"context"
	"fmt"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/internal/repository"
	"github.com/Rin0530/DinnerDecider/backend/internal/service"
)

// recipeUsecase implements the RecipeUsecase interface
type recipeUsecase struct {
	ingredientRepo repository.IngredientRepository
	ollamaService  service.OllamaService
}

// NewRecipeUsecase creates a new instance of RecipeUsecase
func NewRecipeUsecase(
	ingredientRepo repository.IngredientRepository,
	ollamaService service.OllamaService,
) RecipeUsecase {
	return &recipeUsecase{
		ingredientRepo: ingredientRepo,
		ollamaService:  ollamaService,
	}
}

// GetRecipeSuggestion generates recipe suggestions based on available ingredients
func (u *recipeUsecase) GetRecipeSuggestion(ctx context.Context) (*domain.RecipeResponse, error) {
	// Retrieve all ingredients from repository
	ingredients, err := u.ingredientRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients: %w", err)
	}

	// Handle empty ingredients case (len on a nil slice is zero)
	if len(ingredients) == 0 {
		ingredients = []*domain.Ingredient{}
	}

	// Generate recipe suggestions using Ollama service
	recipeResponse, err := u.ollamaService.GenerateRecipeSuggestion(ctx, ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipe suggestion: %w", err)
	}

	return recipeResponse, nil
}
