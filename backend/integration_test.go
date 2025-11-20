package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/internal/handler"
	"github.com/Rin0530/DinnerDecider/backend/internal/repository"
	"github.com/Rin0530/DinnerDecider/backend/internal/service"
	"github.com/Rin0530/DinnerDecider/backend/internal/usecase"
	"github.com/Rin0530/DinnerDecider/backend/pkg/config"
	"github.com/Rin0530/DinnerDecider/backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB creates a MySQL container and returns the database connection
func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()

	// Create MySQL container
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "testpassword",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("3306/tcp"),
			// MySQL prints "ready for connections." (note the trailing dot)
			wait.ForLog("ready for connections."),
		).WithStartupTimeout(120 * time.Second),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}

	// Get container host and port
	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		// attempt to terminate container if possible
		_ = mysqlContainer.Terminate(ctx)
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
		_ = mysqlContainer.Terminate(ctx)
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Waiting for MySQL to be ready is handled by the container wait strategy above.

	// Connect to database
	dbConfig := database.Config{
		Host:         host,
		Port:         port.Int(),
		User:         "root",
		Password:     "testpassword",
		DBName:       "testdb",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		mysqlContainer.Terminate(ctx)
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	schema := `
	CREATE TABLE IF NOT EXISTS ingredients (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity VARCHAR(100),
		expiration_date DATE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_name (name),
		INDEX idx_expiration_date (expiration_date)
	);`

	if _, err := db.Exec(schema); err != nil {
		database.Close(db)
		mysqlContainer.Terminate(ctx)
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		database.Close(db)
		if err := mysqlContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

// setupTestRouter creates a test router with all dependencies
func setupTestRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize dependencies
	ingredientRepo := repository.NewIngredientRepository(db)

	// Mock Ollama service for testing
	timeout, _ := time.ParseDuration("30s")
	ollamaConfig := &config.OllamaConfig{
		Endpoint: "http://localhost:11434",
		Model:    "llama2",
		Timeout:  timeout,
	}
	ollamaService := service.NewOllamaService(ollamaConfig)

	ingredientUsecase := usecase.NewIngredientUsecase(ingredientRepo)
	recipeUsecase := usecase.NewRecipeUsecase(ingredientRepo, ollamaService)

	ingredientHandler := handler.NewIngredientHandler(ingredientUsecase)
	recipeHandler := handler.NewRecipeHandler(recipeUsecase)
	healthHandler := handler.NewHealthHandler(db, ollamaService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", healthHandler.Health)
	router.GET("/health/db", healthHandler.HealthDB)

	api := router.Group("/api")
	{
		ingredients := api.Group("/ingredients")
		{
			ingredients.POST("", ingredientHandler.CreateIngredient)
			ingredients.GET("", ingredientHandler.GetAllIngredients)
			ingredients.GET("/:id", ingredientHandler.GetIngredientByID)
			ingredients.PUT("/:id", ingredientHandler.UpdateIngredient)
			ingredients.DELETE("/:id", ingredientHandler.DeleteIngredient)
		}

		recipes := api.Group("/recipes")
		{
			recipes.POST("/suggestion", recipeHandler.GetRecipeSuggestion)
		}
	}

	return router
}

// TestIngredientCRUDFlow tests the complete CRUD flow for ingredients
func TestIngredientCRUDFlow(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(db)

	// Test 1: Create ingredient
	t.Run("Create Ingredient", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":            "にんじん",
			"quantity":        "2本",
			"expiration_date": "2025-11-01",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Body: %s", w.Body.String())

		var ingredient domain.Ingredient
		err := json.Unmarshal(w.Body.Bytes(), &ingredient)
		assert.NoError(t, err)
		assert.Equal(t, "にんじん", ingredient.Name)
		assert.Equal(t, "2本", ingredient.Quantity)
	})

	// Test 2: Get all ingredients
	t.Run("Get All Ingredients", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/ingredients", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var ingredients []*domain.Ingredient
		if err := json.Unmarshal(w.Body.Bytes(), &ingredients); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(ingredients) == 0 {
			t.Error("Expected at least one ingredient")
		}
	})

	// Test 3: Update ingredient
	t.Run("Update Ingredient", func(t *testing.T) {
		// First create an ingredient to update
		createBody := map[string]interface{}{
			"name":     "豚バラ肉",
			"quantity": "200g",
		}
		body, _ := json.Marshal(createBody)
		req := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var created domain.Ingredient
		err := json.Unmarshal(w.Body.Bytes(), &created)
		assert.NoError(t, err)

		// Update the ingredient
		updateBody := map[string]interface{}{
			"quantity": "300g",
		}
		body, _ = json.Marshal(updateBody)
		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/ingredients/%d", created.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Body: %s", w.Body.String())

		var updated domain.Ingredient
		err = json.Unmarshal(w.Body.Bytes(), &updated)
		assert.NoError(t, err)
		assert.Equal(t, "300g", updated.Quantity)
	})

	// Test 4: Delete ingredient
	t.Run("Delete Ingredient", func(t *testing.T) {
		// First create an ingredient to delete
		createBody := map[string]interface{}{
			"name": "玉ねぎ",
		}
		body, _ := json.Marshal(createBody)
		req := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var created domain.Ingredient
		err := json.Unmarshal(w.Body.Bytes(), &created)
		assert.NoError(t, err)

		// Delete the ingredient
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/ingredients/%d", created.ID), nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify it's deleted
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/ingredients/%d", created.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestValidationErrors tests input validation
func TestValidationErrors(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(db)

	t.Run("Create Ingredient Without Name", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"quantity": "100g",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/ingredients", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update Non-existent Ingredient", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"quantity": "500g",
		}
		body, _ := json.Marshal(updateBody)

		req := httptest.NewRequest(http.MethodPut, "/api/ingredients/99999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Delete Non-existent Ingredient", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/ingredients/99999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestHealthEndpoints tests health check endpoints
func TestHealthEndpoints(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(db)

	t.Run("Health Check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Database Health Check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health/db", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Ollama Health Check", func(t *testing.T) {
		// This test assumes Ollama is not running, so it expects a 503 error.
		// In a real CI/CD environment, you might start an Ollama container
		// and expect a 200 OK.
		req := httptest.NewRequest(http.MethodGet, "/health/ollama", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Since Ollama is not part of the test setup, we expect a service unavailable error.
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

// TestEmptyIngredientsFlow tests behavior with no ingredients
func TestEmptyIngredientsFlow(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupTestRouter(db)

	t.Run("Get All Ingredients When Empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/ingredients", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var ingredients []*domain.Ingredient
		if err := json.Unmarshal(w.Body.Bytes(), &ingredients); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(ingredients) != 0 {
			t.Errorf("Expected empty array, got %d ingredients", len(ingredients))
		}
	})
}
