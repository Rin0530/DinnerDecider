package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return sqlxDB, mock
}

func TestCreate_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	ingredient := &domain.Ingredient{
		Name:         "にんじん",
		Quantity:     "2本",
		PurchaseDate: &purchaseDate,
	}

	mock.ExpectExec("INSERT INTO ingredients").
		WithArgs(ingredient.Name, ingredient.Quantity, ingredient.PurchaseDate, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), ingredient)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), ingredient.ID)
	assert.NotZero(t, ingredient.CreatedAt)
	assert.NotZero(t, ingredient.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	ingredient := &domain.Ingredient{
		Name:     "豚バラ肉",
		Quantity: "200g",
	}

	mock.ExpectExec("INSERT INTO ingredients").
		WithArgs(ingredient.Name, ingredient.Quantity, ingredient.PurchaseDate, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)
	err := repo.Create(context.Background(), ingredient)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create ingredient")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	now := time.Now()
	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "quantity", "purchase_date", "created_at", "updated_at"}).
		AddRow(1, "にんじん", "2本", purchaseDate, now, now).
		AddRow(2, "豚バラ肉", "200g", nil, now, now)

	mock.ExpectQuery("SELECT (.+) FROM ingredients").
		WillReturnRows(rows)

	ingredients, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, ingredients, 2)
	assert.Equal(t, "にんじん", ingredients[0].Name)
	assert.Equal(t, "豚バラ肉", ingredients[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "quantity", "purchase_date", "created_at", "updated_at"})

	mock.ExpectQuery("SELECT (.+) FROM ingredients").
		WillReturnRows(rows)

	ingredients, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, ingredients)
	assert.Len(t, ingredients, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM ingredients").
		WillReturnError(sql.ErrConnDone)

	ingredients, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, ingredients)
	assert.Contains(t, err.Error(), "failed to get all ingredients")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	now := time.Now()
	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "name", "quantity", "purchase_date", "created_at", "updated_at"}).
		AddRow(1, "にんじん", "2本", purchaseDate, now, now)

	mock.ExpectQuery("SELECT (.+) FROM ingredients WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	ingredient, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, ingredient)
	assert.Equal(t, int64(1), ingredient.ID)
	assert.Equal(t, "にんじん", ingredient.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM ingredients WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	ingredient, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, ingredient)
	assert.Contains(t, err.Error(), "ingredient not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM ingredients WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	ingredient, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, ingredient)
	assert.Contains(t, err.Error(), "failed to get ingredient by id")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	purchaseDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	ingredient := &domain.Ingredient{
		ID:           1,
		Name:         "にんじん",
		Quantity:     "3本",
		PurchaseDate: &purchaseDate,
	}

	mock.ExpectExec("UPDATE ingredients").
		WithArgs(ingredient.Name, ingredient.Quantity, ingredient.PurchaseDate, sqlmock.AnyArg(), ingredient.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), ingredient)

	assert.NoError(t, err)
	assert.NotZero(t, ingredient.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	ingredient := &domain.Ingredient{
		ID:       999,
		Name:     "存在しない食材",
		Quantity: "1個",
	}

	mock.ExpectExec("UPDATE ingredients").
		WithArgs(ingredient.Name, ingredient.Quantity, ingredient.PurchaseDate, sqlmock.AnyArg(), ingredient.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), ingredient)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ingredient not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	ingredient := &domain.Ingredient{
		ID:       1,
		Name:     "にんじん",
		Quantity: "3本",
	}

	mock.ExpectExec("UPDATE ingredients").
		WithArgs(ingredient.Name, ingredient.Quantity, ingredient.PurchaseDate, sqlmock.AnyArg(), ingredient.ID).
		WillReturnError(sql.ErrConnDone)

	err := repo.Update(context.Background(), ingredient)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update ingredient")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectExec("DELETE FROM ingredients WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectExec("DELETE FROM ingredients WHERE id = ?").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ingredient not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewIngredientRepository(db)

	mock.ExpectExec("DELETE FROM ingredients WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	err := repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete ingredient")
	assert.NoError(t, mock.ExpectationsWereMet())
}
