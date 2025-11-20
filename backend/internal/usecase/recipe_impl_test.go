package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOllamaService is a mock implementation of OllamaService
type MockOllamaService struct {
	mock.Mock
}

func (m *MockOllamaService) GenerateRecipeSuggestion(ctx context.Context, ingredients []*domain.Ingredient) (*domain.RecipeResponse, error) {
	args := m.Called(ctx, ingredients)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RecipeResponse), args.Error(1)
}

// TestGetRecipeSuggestion_Success tests successful recipe suggestion generation
func TestGetRecipeSuggestion_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	mockService := new(MockOllamaService)
	usecase := NewRecipeUsecase(mockRepo, mockService)

	now := time.Now()
	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	mockIngredients := []*domain.Ingredient{
		{
			ID:           1,
			Name:         "にんじん",
			Quantity:     "2本",
			PurchaseDate: &purchaseDate,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:        2,
			Name:      "豚バラ肉",
			Quantity:  "200g",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockRecipeResponse := &domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "カレーライス",
				Steps:        []string{"野菜を切る", "肉を炒める", "カレールーを入れる"},
				MissingItems: []string{"カレールー"},
			},
			{
				Name:         "豚汁",
				Steps:        []string{"野菜を切る", "豚肉を炒める", "味噌を入れる"},
				MissingItems: []string{"味噌"},
			},
		},
	}

	mockRepo.On("GetAll", mock.Anything).Return(mockIngredients, nil)
	mockService.On("GenerateRecipeSuggestion", mock.Anything, mockIngredients).
		Return(mockRecipeResponse, nil)

	result, err := usecase.GetRecipeSuggestion(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Suggestions, 2)
	assert.Equal(t, "カレーライス", result.Suggestions[0].Name)
	assert.Equal(t, "豚汁", result.Suggestions[1].Name)
	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestGetRecipeSuggestion_EmptyIngredients tests handling of empty ingredient list
func TestGetRecipeSuggestion_EmptyIngredients(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	mockService := new(MockOllamaService)
	usecase := NewRecipeUsecase(mockRepo, mockService)

	mockRecipeResponse := &domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "シンプルな卵かけご飯",
				Steps:        []string{"ご飯を炊く", "卵を割る", "醤油をかける"},
				MissingItems: []string{"卵", "ご飯", "醤油"},
			},
		},
	}

	mockRepo.On("GetAll", mock.Anything).Return([]*domain.Ingredient{}, nil)
	mockService.On("GenerateRecipeSuggestion", mock.Anything, []*domain.Ingredient{}).
		Return(mockRecipeResponse, nil)

	result, err := usecase.GetRecipeSuggestion(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Suggestions, 1)
	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestGetRecipeSuggestion_NilIngredients tests handling of nil ingredient list
func TestGetRecipeSuggestion_NilIngredients(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	mockService := new(MockOllamaService)
	usecase := NewRecipeUsecase(mockRepo, mockService)

	mockRecipeResponse := &domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "シンプルな卵かけご飯",
				Steps:        []string{"ご飯を炊く", "卵を割る", "醤油をかける"},
				MissingItems: []string{"卵", "ご飯", "醤油"},
			},
		},
	}

	mockRepo.On("GetAll", mock.Anything).Return(nil, nil)
	mockService.On("GenerateRecipeSuggestion", mock.Anything, []*domain.Ingredient{}).
		Return(mockRecipeResponse, nil)

	result, err := usecase.GetRecipeSuggestion(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestGetRecipeSuggestion_RepositoryError tests error handling when repository fails
func TestGetRecipeSuggestion_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	mockService := new(MockOllamaService)
	usecase := NewRecipeUsecase(mockRepo, mockService)

	mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))

	result, err := usecase.GetRecipeSuggestion(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get ingredients")
	mockRepo.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GenerateRecipeSuggestion")
}

// TestGetRecipeSuggestion_ServiceError tests error handling when Ollama service fails
func TestGetRecipeSuggestion_ServiceError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	mockService := new(MockOllamaService)
	usecase := NewRecipeUsecase(mockRepo, mockService)

	now := time.Now()
	mockIngredients := []*domain.Ingredient{
		{
			ID:        1,
			Name:      "にんじん",
			Quantity:  "2本",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockRepo.On("GetAll", mock.Anything).Return(mockIngredients, nil)
	mockService.On("GenerateRecipeSuggestion", mock.Anything, mockIngredients).
		Return(nil, errors.New("ollama service unavailable"))

	result, err := usecase.GetRecipeSuggestion(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate recipe suggestion")
	mockRepo.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
