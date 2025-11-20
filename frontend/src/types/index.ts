// 食材の型定義
export interface Ingredient {
  id: number;
  name: string;
  quantity?: string;
  purchase_date?: string; // YYYY-MM-DD形式
  created_at?: string;
  updated_at?: string;
}

// 食材入力用の型定義
export interface IngredientInput {
  name: string;
  quantity?: string;
  purchase_date?: string;
}

// レシピの型定義
export interface Recipe {
  name: string;
  steps: string[];
  missing_items: string[];
}

// レシピ提案レスポンスの型定義
export interface RecipeSuggestionResponse {
  suggestions: Recipe[];
}

// APIエラーの型定義
export interface ApiError {
  error?: string;
  message: string;
  code?: string;
}

// APIレスポンスの型定義
export interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}
