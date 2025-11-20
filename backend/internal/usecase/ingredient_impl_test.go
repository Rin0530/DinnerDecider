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

// MockIngredientRepository is a mock implementation of IngredientRepository
type MockIngredientRepository struct {
	mock.Mock
}

func (m *MockIngredientRepository) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	args := m.Called(ctx, ingredient)
	return args.Error(0)
}

func (m *MockIngredientRepository) GetAll(ctx context.Context) ([]*domain.Ingredient, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientRepository) GetByID(ctx context.Context, id int64) (*domain.Ingredient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ingredient), args.Error(1)
}

func (m *MockIngredientRepository) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	args := m.Called(ctx, ingredient)
	return args.Error(0)
}

func (m *MockIngredientRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestCreateIngredient_Success tests successful ingredient creation
func TestCreateIngredient_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	purchaseDate := "2025-12-31"
	req := CreateIngredientRequest{
		Name:         "にんじん",
		Quantity:     "2本",
		PurchaseDate: &purchaseDate,
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ingredient")).
		Return(nil).
		Run(func(args mock.Arguments) {
			ingredient := args.Get(1).(*domain.Ingredient)
			ingredient.ID = 1
		})

	result, err := usecase.CreateIngredient(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "にんじん", result.Name)
	assert.Equal(t, "2本", result.Quantity)
	assert.NotNil(t, result.PurchaseDate)
	mockRepo.AssertExpectations(t)
}

// TestCreateIngredient_MissingName tests validation error when name is missing
func TestCreateIngredient_MissingName(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	req := CreateIngredientRequest{
		Name:     "",
		Quantity: "2本",
	}

	result, err := usecase.CreateIngredient(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateIngredient_InvalidPurchaseDate tests error handling for invalid date format
func TestCreateIngredient_InvalidPurchaseDate(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	invalidDate := "invalid-date"
	req := CreateIngredientRequest{
		Name:         "にんじん",
		Quantity:     "2本",
		PurchaseDate: &invalidDate,
	}

	result, err := usecase.CreateIngredient(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid purchase_date format")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateIngredient_RepositoryError tests error handling when repository fails
func TestCreateIngredient_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	req := CreateIngredientRequest{
		Name:     "にんじん",
		Quantity: "2本",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Ingredient")).
		Return(errors.New("database error"))

	result, err := usecase.CreateIngredient(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create ingredient")
	mockRepo.AssertExpectations(t)
}

// TestGetAllIngredients_Success tests successful retrieval of all ingredients
func TestGetAllIngredients_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

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

	mockRepo.On("GetAll", mock.Anything).Return(mockIngredients, nil)

	result, err := usecase.GetAllIngredients(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "にんじん", result[0].Name)
	assert.Equal(t, "豚バラ肉", result[1].Name)
	mockRepo.AssertExpectations(t)
}

// TestGetAllIngredients_EmptyResult tests handling of empty ingredient list
func TestGetAllIngredients_EmptyResult(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	mockRepo.On("GetAll", mock.Anything).Return(nil, nil)

	result, err := usecase.GetAllIngredients(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

// TestGetAllIngredients_RepositoryError tests error handling when repository fails
func TestGetAllIngredients_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("database error"))

	result, err := usecase.GetAllIngredients(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get all ingredients")
	mockRepo.AssertExpectations(t)
}

// TestGetIngredientByID_Success tests successful ingredient retrieval by ID
func TestGetIngredientByID_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	expectedIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedIngredient, nil)

	result, err := usecase.GetIngredientByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "にんじん", result.Name)
	mockRepo.AssertExpectations(t)
}

// TestGetIngredientByID_NotFound tests error handling when ingredient doesn't exist
func TestGetIngredientByID_NotFound(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	result, err := usecase.GetIngredientByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get ingredient by id")
	mockRepo.AssertExpectations(t)
}

// TestGetIngredientByID_RepositoryError tests error handling when repository fails
func TestGetIngredientByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))

	result, err := usecase.GetIngredientByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get ingredient by id")
	mockRepo.AssertExpectations(t)
}

// TestUpdateIngredient_Success tests successful ingredient update
func TestUpdateIngredient_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newName := "大根"
	newQuantity := "3本"
	req := UpdateIngredientRequest{
		Name:     &newName,
		Quantity: &newQuantity,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Ingredient")).Return(nil)

	result, err := usecase.UpdateIngredient(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "大根", result.Name)
	assert.Equal(t, "3本", result.Quantity)
	mockRepo.AssertExpectations(t)
}

// TestUpdateIngredient_PartialUpdate tests partial update of ingredient fields
func TestUpdateIngredient_PartialUpdate(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newQuantity := "3本"
	req := UpdateIngredientRequest{
		Quantity: &newQuantity,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Ingredient")).Return(nil)

	result, err := usecase.UpdateIngredient(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "にんじん", result.Name)
	assert.Equal(t, "3本", result.Quantity)
	mockRepo.AssertExpectations(t)
}

// TestUpdateIngredient_NotFound tests error handling when ingredient doesn't exist
func TestUpdateIngredient_NotFound(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	newName := "大根"
	req := UpdateIngredientRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("ingredient not found"))

	result, err := usecase.UpdateIngredient(context.Background(), 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get ingredient")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

// TestUpdateIngredient_InvalidPurchaseDate tests error handling for invalid date format
func TestUpdateIngredient_InvalidPurchaseDate(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	invalidDate := "invalid-date"
	req := UpdateIngredientRequest{
		PurchaseDate: &invalidDate,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)

	result, err := usecase.UpdateIngredient(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid purchase_date format")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

// TestUpdateIngredient_RepositoryError tests error handling when repository fails
func TestUpdateIngredient_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newName := "大根"
	req := UpdateIngredientRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Ingredient")).
		Return(errors.New("database error"))

	result, err := usecase.UpdateIngredient(context.Background(), 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update ingredient")
	mockRepo.AssertExpectations(t)
}

// TestDeleteIngredient_Success tests successful ingredient deletion
func TestDeleteIngredient_Success(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := usecase.DeleteIngredient(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteIngredient_NotFound tests error handling when ingredient doesn't exist
func TestDeleteIngredient_NotFound(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("ingredient not found"))

	err := usecase.DeleteIngredient(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get ingredient")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

// TestDeleteIngredient_RepositoryError tests error handling when repository fails
func TestDeleteIngredient_RepositoryError(t *testing.T) {
	mockRepo := new(MockIngredientRepository)
	usecase := NewIngredientUsecase(mockRepo)

	now := time.Now()
	existingIngredient := &domain.Ingredient{
		ID:        1,
		Name:      "にんじん",
		Quantity:  "2本",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingIngredient, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	err := usecase.DeleteIngredient(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete ingredient")
	mockRepo.AssertExpectations(t)
}
