package handler

import (
	"net/http"
	"strings"

	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// RecipeHandler handles HTTP requests for recipe operations
type RecipeHandler struct {
	recipeUsecase usecase.RecipeUsecase
}

// NewRecipeHandler creates a new RecipeHandler instance
func NewRecipeHandler(recipeUsecase usecase.RecipeUsecase) *RecipeHandler {
	return &RecipeHandler{
		recipeUsecase: recipeUsecase,
	}
}

// GetRecipeSuggestion handles POST /recipes/suggestion
// @Summary 献立提案を取得
// @Description 登録されている食材を基にAIが献立を提案します
// @Tags recipes
// @Accept json
// @Produce json
// @Success 200 {object} domain.RecipeResponse "献立提案のリスト"
// @Failure 500 {object} usecase.ErrorResponse "内部サーバーエラー"
// @Failure 503 {object} usecase.ErrorResponse "サービス利用不可（AI APIが利用できない場合）"
// @Router /recipes/suggestion [post]
func (h *RecipeHandler) GetRecipeSuggestion(c *gin.Context) {
	// Call usecase to get recipe suggestions
	recipeResponse, err := h.recipeUsecase.GetRecipeSuggestion(c.Request.Context())
	if err != nil {
		// Check if it's a service unavailability error (Ollama API)
		if strings.Contains(err.Error(), "ollama") || 
		   strings.Contains(err.Error(), "timeout") ||
		   strings.Contains(err.Error(), "connection") {
			respondServiceUnavailable(c, "Recipe suggestion service is currently unavailable")
			return
		}
		
		// Handle other errors
		handleError(c, err)
		return
	}

	// Return recipe suggestions
	c.JSON(http.StatusOK, recipeResponse)
}
