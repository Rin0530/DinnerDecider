package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIngredientUsecase is a mock implementation of IngredientUsecase
type MockIngredientUsecase struct {
	mock.Mock
}

func (m *MockIngredientUsecase) CreateIngredient(ctx context.Context, req usecase.CreateIngredientRequest) (*domain.Ingredient, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientUsecase) GetAllIngredients(ctx context.Context) ([]*domain.Ingredient, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientUsecase) GetIngredientByID(ctx context.Context, id int64) (*domain.Ingredient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientUsecase) UpdateIngredient(ctx context.Context, id int64, req usecase.UpdateIngredientRequest) (*domain.Ingredient, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientUsecase) DeleteIngredient(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestCreateIngredient_Success tests successful ingredient creation
func TestCreateIngredient_Success(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/ingredients", handler.CreateIngredient)

	now := time.Now()
	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	expectedIngredient := &domain.Ingredient{
		ID:           1,
		Name:         "にんじん",
		Quantity:     "2本",
		PurchaseDate: &purchaseDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	purchaseDateStr := "2025-12-31"
	reqBody := usecase.CreateIngredientRequest{
		Name:         "にんじん",
		Quantity:     "2本",
		PurchaseDate: &purchaseDateStr,
	}

	mockUsecase.On("CreateIngredient", mock.Anything, reqBody).Return(expectedIngredient, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "にんじん", response.Name)
	assert.Equal(t, "2本", response.Quantity)
	mockUsecase.AssertExpectations(t)
}

// TestCreateIngredient_MissingName tests validation error when name is missing
func TestCreateIngredient_MissingName(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/ingredients", handler.CreateIngredient)

	reqBody := map[string]interface{}{
		"quantity": "2本",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	mockUsecase.AssertNotCalled(t, "CreateIngredient")
}

// TestCreateIngredient_InvalidJSON tests error handling for invalid JSON
func TestCreateIngredient_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/ingredients", handler.CreateIngredient)

	req := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	mockUsecase.AssertNotCalled(t, "CreateIngredient")
}

// TestCreateIngredient_UsecaseError tests error handling when usecase fails
func TestCreateIngredient_UsecaseError(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/ingredients", handler.CreateIngredient)

	reqBody := usecase.CreateIngredientRequest{
		Name:     "にんじん",
		Quantity: "2本",
	}

	mockUsecase.On("CreateIngredient", mock.Anything, reqBody).
		Return(nil, errors.New("database error"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/ingredients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestGetIngredientByID_Success tests successful ingredient retrieval by ID
func TestGetIngredientByID_Success(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients/:id", handler.GetIngredientByID)

	now := time.Now()
	expectedIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUsecase.On("GetIngredientByID", mock.Anything, int64(1)).Return(expectedIngredient, nil)

	req := httptest.NewRequest(http.MethodGet, "/ingredients/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "にんじん", response.Name)
	mockUsecase.AssertExpectations(t)
}

// TestGetIngredientByID_InvalidID tests error handling for invalid ID
func TestGetIngredientByID_InvalidID(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients/:id", handler.GetIngredientByID)

	req := httptest.NewRequest(http.MethodGet, "/ingredients/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	assert.Contains(t, response.Message, "Invalid ingredient ID")
	mockUsecase.AssertNotCalled(t, "GetIngredientByID")
}

// TestGetIngredientByID_NotFound tests error handling when ingredient doesn't exist
func TestGetIngredientByID_NotFound(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients/:id", handler.GetIngredientByID)

	mockUsecase.On("GetIngredientByID", mock.Anything, int64(999)).Return(nil, sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/ingredients/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestGetAllIngredients_Success tests successful retrieval of all ingredients
func TestGetAllIngredients_Success(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients", handler.GetAllIngredients)

	now := time.Now()
	mockIngredients := []*domain.Ingredient{
		{
			ID:        1,
			Name:      "にんじん",
			Quantity:  "2本",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "豚バラ肉",
			Quantity:  "200g",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockUsecase.On("GetAllIngredients", mock.Anything).Return(mockIngredients, nil)

	req := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*domain.Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "にんじん", response[0].Name)
	assert.Equal(t, "豚バラ肉", response[1].Name)
	mockUsecase.AssertExpectations(t)
}

// TestGetAllIngredients_EmptyResult tests handling of empty ingredient list
func TestGetAllIngredients_EmptyResult(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients", handler.GetAllIngredients)

	mockUsecase.On("GetAllIngredients", mock.Anything).Return([]*domain.Ingredient{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*domain.Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
	mockUsecase.AssertExpectations(t)
}

// TestGetAllIngredients_UsecaseError tests error handling when usecase fails
func TestGetAllIngredients_UsecaseError(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.GET("/ingredients", handler.GetAllIngredients)

	mockUsecase.On("GetAllIngredients", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/ingredients", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestUpdateIngredient_Success tests successful ingredient update
func TestUpdateIngredient_Success(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.PUT("/ingredients/:id", handler.UpdateIngredient)

	now := time.Now()
	updatedIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "大根",
		Quantity:  "3本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newName := "大根"
	newQuantity := "3本"
	reqBody := usecase.UpdateIngredientRequest{
		Name:     &newName,
		Quantity: &newQuantity,
	}

	mockUsecase.On("UpdateIngredient", mock.Anything, int64(1), reqBody).
		Return(updatedIngredient, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/ingredients/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "大根", response.Name)
	assert.Equal(t, "3本", response.Quantity)
	mockUsecase.AssertExpectations(t)
}

// TestUpdateIngredient_InvalidID tests error handling for invalid ID
func TestUpdateIngredient_InvalidID(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.PUT("/ingredients/:id", handler.UpdateIngredient)

	newName := "大根"
	reqBody := usecase.UpdateIngredientRequest{
		Name: &newName,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/ingredients/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	assert.Contains(t, response.Message, "Invalid ingredient ID")
	mockUsecase.AssertNotCalled(t, "UpdateIngredient")
}

// TestUpdateIngredient_NotFound tests error handling when ingredient doesn't exist
func TestUpdateIngredient_NotFound(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.PUT("/ingredients/:id", handler.UpdateIngredient)

	newName := "大根"
	reqBody := usecase.UpdateIngredientRequest{
		Name: &newName,
	}

	mockUsecase.On("UpdateIngredient", mock.Anything, int64(999), reqBody).
		Return(nil, sql.ErrNoRows)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/ingredients/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)
	mockUsecase.AssertExpectations(t)
}

// TestUpdateIngredient_InvalidJSON tests error handling for invalid JSON
func TestUpdateIngredient_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.PUT("/ingredients/:id", handler.UpdateIngredient)

	req := httptest.NewRequest(http.MethodPut, "/ingredients/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	mockUsecase.AssertNotCalled(t, "UpdateIngredient")
}

// TestDeleteIngredient_Success tests successful ingredient deletion
func TestDeleteIngredient_Success(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.DELETE("/ingredients/:id", handler.DeleteIngredient)

	mockUsecase.On("DeleteIngredient", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/ingredients/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	mockUsecase.AssertExpectations(t)
}

// TestDeleteIngredient_InvalidID tests error handling for invalid ID
func TestDeleteIngredient_InvalidID(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.DELETE("/ingredients/:id", handler.DeleteIngredient)

	req := httptest.NewRequest(http.MethodDelete, "/ingredients/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	assert.Contains(t, response.Message, "Invalid ingredient ID")
	mockUsecase.AssertNotCalled(t, "DeleteIngredient")
}

// TestDeleteIngredient_NotFound tests error handling when ingredient doesn't exist
func TestDeleteIngredient_NotFound(t *testing.T) {
	mockUsecase := new(MockIngredientUsecase)
	handler := NewIngredientHandler(mockUsecase)
	router := setupTestRouter()
	router.DELETE("/ingredients/:id", handler.DeleteIngredient)

	mockUsecase.On("DeleteIngredient", mock.Anything, int64(999)).Return(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodDelete, "/ingredients/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response usecase.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", response.Error)
	mockUsecase.AssertExpectations(t)
}
