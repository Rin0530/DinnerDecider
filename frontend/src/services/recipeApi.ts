import type { RecipeSuggestionResponse } from '../types';
import { apiClient } from './api';

/**
 * レシピAPI
 * レシピの提案をリクエストするAPIメソッド
 */
export const recipeApi = {
  /**
   * 献立の提案をリクエスト
   */
  requestSuggestion: (): Promise<RecipeSuggestionResponse> => {
    return apiClient.post<RecipeSuggestionResponse>('/recipes/suggestion', {});
  },
};
