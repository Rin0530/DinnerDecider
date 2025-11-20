package handler

import (
	"net/http"
	"strconv"

	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// IngredientHandler handles HTTP requests for ingredient operations
type IngredientHandler struct {
	ingredientUsecase usecase.IngredientUsecase
}

// NewIngredientHandler creates a new IngredientHandler instance
func NewIngredientHandler(ingredientUsecase usecase.IngredientUsecase) *IngredientHandler {
	return &IngredientHandler{
		ingredientUsecase: ingredientUsecase,
	}
}

// @Summary      食材を作成
// @Description  新しい食材を冷蔵庫に追加します。
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param        ingredient body usecase.CreateIngredientRequest true "作成する食材の情報"
// @Success      201 {object} domain.Ingredient "作成された食材"
// @Failure      400 {object} usecase.ErrorResponse "リクエストが不正です"
// @Failure      500 {object} usecase.ErrorResponse "サーバー内部エラー"
// @Router       /ingredients [post]
// CreateIngredient handles POST /ingredients
func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var req usecase.CreateIngredientRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	// Call usecase
	ingredient, err := h.ingredientUsecase.CreateIngredient(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Return created ingredient with 201 status
	c.JSON(http.StatusCreated, ingredient)
}

// @Summary      すべての食材を取得
// @Description  冷蔵庫にあるすべての食材のリストを取得します。
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Success      200 {array} domain.Ingredient "食材のリスト"
// @Failure      500 {object} usecase.ErrorResponse "サーバー内部エラー"
// @Router       /ingredients [get]
// GetAllIngredients handles GET /ingredients
func (h *IngredientHandler) GetAllIngredients(c *gin.Context) {
	// Call usecase
	ingredients, err := h.ingredientUsecase.GetAllIngredients(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	// Return ingredients array (empty array if no ingredients)
	c.JSON(http.StatusOK, ingredients)
}

// @Summary      IDで食材を取得
// @Description  指定されたIDの食材情報を取得します。
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "食材ID"
// @Success      200 {object} domain.Ingredient "食材の情報"
// @Failure      400 {object} usecase.ErrorResponse "リクエストが不正です"
// @Failure      404 {object} usecase.ErrorResponse "食材が見つかりません"
// @Failure      500 {object} usecase.ErrorResponse "サーバー内部エラー"
// @Router       /ingredients/{id} [get]
// GetIngredientByID handles GET /ingredients/:id
func (h *IngredientHandler) GetIngredientByID(c *gin.Context) {
	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondBadRequest(c, "Invalid ingredient ID")
		return
	}

	// Call usecase
	ingredient, err := h.ingredientUsecase.GetIngredientByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

// @Summary      食材を更新
// @Description  指定されたIDの食材情報を更新します。
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "食材ID"
// @Param        ingredient body usecase.UpdateIngredientRequest true "更新する食材の情報"
// @Success      200 {object} domain.Ingredient "更新された食材"
// @Failure      400 {object} usecase.ErrorResponse "リクエストが不正です"
// @Failure      404 {object} usecase.ErrorResponse "食材が見つかりません"
// @Failure      500 {object} usecase.ErrorResponse "サーバー内部エラー"
// @Router       /ingredients/{id} [put]
// UpdateIngredient handles PUT /ingredients/:id
func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondBadRequest(c, "Invalid ingredient ID")
		return
	}

	var req usecase.UpdateIngredientRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	// Call usecase
	ingredient, err := h.ingredientUsecase.UpdateIngredient(c.Request.Context(), id, req)
	if err != nil {
		handleError(c, err)
		return
	}

	// Return updated ingredient
	c.JSON(http.StatusOK, ingredient)
}

// @Summary      食材を削除
// @Description  指定されたIDの食材を削除します。
// @Tags         ingredients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "食材ID"
// @Success      204 "削除成功"
// @Failure      400 {object} usecase.ErrorResponse "リクエストが不正です"
// @Failure      500 {object} usecase.ErrorResponse "サーバー内部エラー"
// @Router       /ingredients/{id} [delete]
// DeleteIngredient handles DELETE /ingredients/:id
func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondBadRequest(c, "Invalid ingredient ID")
		return
	}

	// Call usecase
	err = h.ingredientUsecase.DeleteIngredient(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	// Return 204 No Content on successful deletion
	c.Status(http.StatusNoContent)
}
