package usecase

// CreateIngredientRequest represents the request body for creating a new ingredient
type CreateIngredientRequest struct {
	Name         string  `json:"name" binding:"required"`
	Quantity     string  `json:"quantity"`
	PurchaseDate *string `json:"purchase_date"` // YYYY-MM-DD format
}

// UpdateIngredientRequest represents the request body for updating an ingredient
type UpdateIngredientRequest struct {
	Name         *string `json:"name"`
	Quantity     *string `json:"quantity"`
	PurchaseDate *string `json:"purchase_date"` // YYYY-MM-DD format
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
