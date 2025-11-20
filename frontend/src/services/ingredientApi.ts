import type { Ingredient, IngredientInput } from '../types';
import { apiClient } from './api';

/**
 * 食材API
 * 食材の取得、作成、更新、削除を行うAPIメソッド
 */
export const ingredientApi = {
  /**
   * 全ての食材を取得
   */
  getAll: (): Promise<Ingredient[]> => {
    return apiClient.get<Ingredient[]>('/ingredients');
  },

  /**
   * 新しい食材を作成
   */
  create: (data: IngredientInput): Promise<Ingredient> => {
    return apiClient.post<Ingredient>('/ingredients', data);
  },

  /**
   * 食材を更新
   */
  update: (id: number, data: IngredientInput): Promise<Ingredient> => {
    return apiClient.put<Ingredient>(`/ingredients/${id}`, data);
  },

  /**
   * 食材を削除
   */
  delete: (id: number): Promise<void> => {
    return apiClient.delete<void>(`/ingredients/${id}`);
  },
};
