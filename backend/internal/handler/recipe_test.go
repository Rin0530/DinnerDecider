package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecipeUsecase is a mock implementation of RecipeUsecase
type MockRecipeUsecase struct {
	mock.Mock
}

func (m *MockRecipeUsecase) GetRecipeSuggestion(ctx context.Context) (*domain.RecipeResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RecipeResponse), args.Error(1)
}

// TestGetRecipeSuggestion_Success tests successful recipe suggestion retrieval
func TestGetRecipeSuggestion_Success(t *testing.T) {
	mockUsecase := new(MockRecipeUsecase)
	handler := NewRecipeHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/recipes/suggestion", handler.GetRecipeSuggestion)

	mockResponse := &domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "野菜炒め",
				Steps:        []string{"野菜を切る", "フライパンで炒める", "調味料で味付けする"},
				MissingItems: []string{},
			},
			{
				Name:         "豚肉の生姜焼き",
				Steps:        []string{"豚肉を切る", "生姜を準備する", "フライパンで焼く"},
				MissingItems: []string{"生姜"},
			},
		},
	}

	mockUsecase.On("GetRecipeSuggestion", mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodPost, "/recipes/suggestion", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response domain.RecipeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Suggestions, 2)
	assert.Equal(t, "野菜炒め", response.Suggestions[0].Name)
	assert.Equal(t, "豚肉の生姜焼き", response.Suggestions[1].Name)
	mockUsecase.AssertExpectations(t)
}

// TestGetRecipeSuggestion_OllamaServiceUnavailable tests error handling when Ollama service is unavailable
func TestGetRecipeSuggestion_OllamaServiceUnavailable(t *testing.T) {
	mockUsecase := new(MockRecipeUsecase)
	handler := NewRecipeHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/recipes/suggestion", handler.GetRecipeSuggestion)

	mockUsecase.On("GetRecipeSuggestion", mock.Anything).
		Return(nil, errors.New("ollama service connection failed"))

	req := httptest.NewRequest(http.MethodPost, "/recipes/suggestion", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "service_unavailable", response.Error)
	assert.Contains(t, response.Message, "Recipe suggestion service is currently unavailable")
	mockUsecase.AssertExpectations(t)
}

// TestGetRecipeSuggestion_TimeoutError tests error handling when request times out
func TestGetRecipeSuggestion_TimeoutError(t *testing.T) {
	mockUsecase := new(MockRecipeUsecase)
	handler := NewRecipeHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/recipes/suggestion", handler.GetRecipeSuggestion)

	mockUsecase.On("GetRecipeSuggestion", mock.Anything).
		Return(nil, errors.New("request timeout"))

	req := httptest.NewRequest(http.MethodPost, "/recipes/suggestion", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "service_unavailable", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestGetRecipeSuggestion_ConnectionError tests error handling when connection fails
func TestGetRecipeSuggestion_ConnectionError(t *testing.T) {
	mockUsecase := new(MockRecipeUsecase)
	handler := NewRecipeHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/recipes/suggestion", handler.GetRecipeSuggestion)

	mockUsecase.On("GetRecipeSuggestion", mock.Anything).
		Return(nil, errors.New("connection refused"))

	req := httptest.NewRequest(http.MethodPost, "/recipes/suggestion", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	
	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "service_unavailable", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestGetRecipeSuggestion_InternalError tests error handling for other internal errors
func TestGetRecipeSuggestion_InternalError(t *testing.T) {
	mockUsecase := new(MockRecipeUsecase)
	handler := NewRecipeHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/recipes/suggestion", handler.GetRecipeSuggestion)

	mockUsecase.On("GetRecipeSuggestion", mock.Anything).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodPost, "/recipes/suggestion", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	mockUsecase.AssertExpectations(t)
}
