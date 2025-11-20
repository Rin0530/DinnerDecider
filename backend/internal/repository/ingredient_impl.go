package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// ingredientRepository is the MySQL implementation of IngredientRepository
type ingredientRepository struct {
	db *sqlx.DB
}

// NewIngredientRepository creates a new instance of IngredientRepository
func NewIngredientRepository(db *sqlx.DB) IngredientRepository {
	return &ingredientRepository{
		db: db,
	}
}

// Create inserts a new ingredient into the database
func (r *ingredientRepository) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	query := `
		INSERT INTO ingredients (name, quantity, purchase_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	ingredient.CreatedAt = now
	ingredient.UpdatedAt = now

	result, err := r.db.ExecContext(
		ctx,
		query,
		ingredient.Name,
		ingredient.Quantity,
		ingredient.PurchaseDate,
		ingredient.CreatedAt,
		ingredient.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create ingredient: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	ingredient.ID = id
	return nil
}

// GetAll retrieves all ingredients from the database
func (r *ingredientRepository) GetAll(ctx context.Context) ([]*domain.Ingredient, error) {
	query := `
		SELECT id, name, quantity, purchase_date, created_at, updated_at
		FROM ingredients
		ORDER BY created_at DESC
	`

	var ingredients []*domain.Ingredient
	err := r.db.SelectContext(ctx, &ingredients, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ingredients: %w", err)
	}

	// Return empty slice instead of nil if no ingredients found
	if ingredients == nil {
		ingredients = []*domain.Ingredient{}
	}

	return ingredients, nil
}

// GetByID retrieves a single ingredient by its ID
func (r *ingredientRepository) GetByID(ctx context.Context, id int64) (*domain.Ingredient, error) {
	query := `
		SELECT id, name, quantity, purchase_date, created_at, updated_at
		FROM ingredients
		WHERE id = ?
	`

	var ingredient domain.Ingredient
	err := r.db.GetContext(ctx, &ingredient, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ingredient not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get ingredient by id: %w", err)
	}

	return &ingredient, nil
}

// Update modifies an existing ingredient in the database
func (r *ingredientRepository) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	query := `
		UPDATE ingredients
		SET name = ?, quantity = ?, purchase_date = ?, updated_at = ?
		WHERE id = ?
	`

	ingredient.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		ingredient.Name,
		ingredient.Quantity,
		ingredient.PurchaseDate,
		ingredient.UpdatedAt,
		ingredient.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update ingredient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ingredient not found")
	}

	return nil
}

// Delete removes an ingredient from the database by its ID
func (r *ingredientRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM ingredients
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ingredient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ingredient not found")
	}

	return nil
}
