package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Rin0530/DinnerDecider/backend/internal/domain"
	"github.com/Rin0530/DinnerDecider/backend/pkg/config"
)

// ollamaServiceImpl implements OllamaService interface
type ollamaServiceImpl struct {
	config     *config.OllamaConfig
	httpClient *http.Client
}

// NewOllamaService creates a new instance of OllamaService
func NewOllamaService(cfg *config.OllamaConfig) OllamaService {
	return &ollamaServiceImpl{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// ollamaRequest represents the request structure for Ollama API
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

// ollamaResponse represents the response structure from Ollama API
type ollamaResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// promptTemplate is the template for generating recipe suggestions
const promptTemplate = `あなたはプロの料理人兼管理栄養士です。以下の食材を使って作れる、美味しくて簡単な夕食の献立を3つ提案してください。
それぞれの献立には、料理名、簡単な作り方、そして不足している食材（もしあれば）を記載してください。
回答は必ずJSON形式で、以下のフォーマットに従ってください。

{
  "suggestions": [
    {
      "name": "料理名",
      "steps": ["手順1", "手順2", "手順3"],
      "missing_items": ["不足している食材1"]
    }
  ]
}

# 利用可能な食材
%s`

// GenerateRecipeSuggestion generates recipe suggestions based on available ingredients
func (s *ollamaServiceImpl) GenerateRecipeSuggestion(ctx context.Context, ingredients []*domain.Ingredient) (*domain.RecipeResponse, error) {
	// Format ingredients list
	ingredientsList := s.formatIngredients(ingredients)

	// Build prompt
	prompt := fmt.Sprintf(promptTemplate, ingredientsList)

	// Create request payload
	reqPayload := ollamaRequest{
		Model:  s.config.Model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	endpoint := fmt.Sprintf("%s/api/generate", s.config.Endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse Ollama response
	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Ollama response: %w", err)
	}

	// Parse the recipe response from the LLM output
	var recipeResp domain.RecipeResponse
	if err := json.Unmarshal([]byte(ollamaResp.Response), &recipeResp); err != nil {
		return nil, fmt.Errorf("failed to parse recipe response: %w", err)
	}

	return &recipeResp, nil
}

// formatIngredients formats the ingredients list into a human-readable string
func (s *ollamaServiceImpl) formatIngredients(ingredients []*domain.Ingredient) string {
	if len(ingredients) == 0 {
		return "食材がありません"
	}

	var parts []string
	for _, ing := range ingredients {
		if ing.Quantity != "" {
			parts = append(parts, fmt.Sprintf("%s(%s)", ing.Name, ing.Quantity))
		} else {
			parts = append(parts, ing.Name)
		}
	}

	return strings.Join(parts, ", ")
}
