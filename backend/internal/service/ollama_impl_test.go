package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/pkg/config"
)

func TestGenerateRecipeSuggestion_Success(t *testing.T) {
	// Create mock recipe response
	mockRecipeResponse := domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "カレーライス",
				Steps:        []string{"野菜を切る", "肉を炒める", "カレールーを入れる"},
				MissingItems: []string{"カレールー"},
			},
		},
	}

	// Create mock Ollama response
	mockOllamaResponse := ollamaResponse{
		Model:     "llama2",
		CreatedAt: time.Now(),
		Response:  mustMarshalJSON(mockRecipeResponse),
		Done:      true,
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockOllamaResponse)
	}))
	defer server.Close()

	// Create service with mock server URL
	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	// Create test ingredients
	ingredients := []*domain.Ingredient{
		{
			ID:       1,
			Name:     "にんじん",
			Quantity: "2本",
		},
		{
			ID:       2,
			Name:     "豚バラ肉",
			Quantity: "200g",
		},
	}

	// Execute
	ctx := context.Background()
	result, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Suggestions) != 1 {
		t.Errorf("Expected 1 suggestion, got %d", len(result.Suggestions))
	}

	if result.Suggestions[0].Name != "カレーライス" {
		t.Errorf("Expected recipe name 'カレーライス', got '%s'", result.Suggestions[0].Name)
	}
}

func TestGenerateRecipeSuggestion_EmptyIngredients(t *testing.T) {
	// Create mock recipe response for empty ingredients
	mockRecipeResponse := domain.RecipeResponse{
		Suggestions: []domain.RecipeSuggestion{
			{
				Name:         "シンプルな卵かけご飯",
				Steps:        []string{"ご飯を炊く", "卵を割る", "醤油をかける"},
				MissingItems: []string{"卵", "ご飯", "醤油"},
			},
		},
	}

	mockOllamaResponse := ollamaResponse{
		Model:     "llama2",
		CreatedAt: time.Now(),
		Response:  mustMarshalJSON(mockRecipeResponse),
		Done:      true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockOllamaResponse)
	}))
	defer server.Close()

	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	// Execute with empty ingredients
	ctx := context.Background()
	result, err := service.GenerateRecipeSuggestion(ctx, []*domain.Ingredient{})

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGenerateRecipeSuggestion_Timeout(t *testing.T) {
	// Create server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create service with short timeout
	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  50 * time.Millisecond,
	}
	service := NewOllamaService(cfg)

	ingredients := []*domain.Ingredient{
		{Name: "にんじん", Quantity: "2本"},
	}

	// Execute
	ctx := context.Background()
	_, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify timeout error
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
}

func TestGenerateRecipeSuggestion_HTTPError(t *testing.T) {
	// Create server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	ingredients := []*domain.Ingredient{
		{Name: "にんじん", Quantity: "2本"},
	}

	// Execute
	ctx := context.Background()
	_, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGenerateRecipeSuggestion_InvalidOllamaResponse(t *testing.T) {
	// Create server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	ingredients := []*domain.Ingredient{
		{Name: "にんじん", Quantity: "2本"},
	}

	// Execute
	ctx := context.Background()
	_, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify error
	if err == nil {
		t.Fatal("Expected JSON unmarshal error, got nil")
	}
}

func TestGenerateRecipeSuggestion_InvalidRecipeResponse(t *testing.T) {
	// Create Ollama response with invalid recipe JSON
	mockOllamaResponse := ollamaResponse{
		Model:     "llama2",
		CreatedAt: time.Now(),
		Response:  "invalid recipe json",
		Done:      true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockOllamaResponse)
	}))
	defer server.Close()

	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	ingredients := []*domain.Ingredient{
		{Name: "にんじん", Quantity: "2本"},
	}

	// Execute
	ctx := context.Background()
	_, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify error
	if err == nil {
		t.Fatal("Expected recipe parse error, got nil")
	}
}

func TestGenerateRecipeSuggestion_ContextCancellation(t *testing.T) {
	// Create server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.OllamaConfig{
		Endpoint: server.URL,
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg)

	ingredients := []*domain.Ingredient{
		{Name: "にんじん", Quantity: "2本"},
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Execute
	_, err := service.GenerateRecipeSuggestion(ctx, ingredients)

	// Verify context cancellation error
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}
}

func TestFormatIngredients(t *testing.T) {
	cfg := &config.OllamaConfig{
		Endpoint: "http://localhost:11434",
		Model:    "llama2",
		Timeout:  30 * time.Second,
	}
	service := NewOllamaService(cfg).(*ollamaServiceImpl)

	tests := []struct {
		name        string
		ingredients []*domain.Ingredient
		expected    string
	}{
		{
			name:        "Empty ingredients",
			ingredients: []*domain.Ingredient{},
			expected:    "食材がありません",
		},
		{
			name: "Single ingredient with quantity",
			ingredients: []*domain.Ingredient{
				{Name: "にんじん", Quantity: "2本"},
			},
			expected: "にんじん(2本)",
		},
		{
			name: "Single ingredient without quantity",
			ingredients: []*domain.Ingredient{
				{Name: "にんじん", Quantity: ""},
			},
			expected: "にんじん",
		},
		{
			name: "Multiple ingredients",
			ingredients: []*domain.Ingredient{
				{Name: "にんじん", Quantity: "2本"},
				{Name: "豚バラ肉", Quantity: "200g"},
				{Name: "玉ねぎ", Quantity: ""},
			},
			expected: "にんじん(2本), 豚バラ肉(200g), 玉ねぎ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.formatIngredients(tt.ingredients)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Helper function to marshal JSON or panic
func mustMarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
