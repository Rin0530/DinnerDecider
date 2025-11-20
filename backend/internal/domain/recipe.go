package domain

// RecipeSuggestion represents a recipe suggestion from LLM
type RecipeSuggestion struct {
	Name         string   `json:"name"`
	Steps        []string `json:"steps"`
	MissingItems []string `json:"missing_items"`
}

// RecipeResponse represents the response containing multiple suggestions
type RecipeResponse struct {
	Suggestions []RecipeSuggestion `json:"suggestions"`
}
