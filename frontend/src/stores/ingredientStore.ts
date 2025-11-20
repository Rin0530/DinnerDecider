import { create } from 'zustand';
import type { Ingredient } from '../types';

// 食材ストアの状態とアクションの型定義
interface IngredientState {
  ingredients: Ingredient[]; // 食材リスト
  isLoading: boolean; // ローディング状態
  error: string | null; // エラーメッセージ
  setIngredients: (ingredients: Ingredient[]) => void; // 食材リストを設定
  addIngredient: (ingredient: Ingredient) => void; // 食材を追加
  updateIngredient: (id: number, ingredient: Ingredient) => void; // 食材を更新
  removeIngredient: (id: number) => void; // 食材を削除
  setLoading: (isLoading: boolean) => void; // ローディング状態を設定
  setError: (error: string | null) => void; // エラーを設定
}

/**
 * 食材ストア
 * Zustandを使用した食材の状態管理
 */
export const useIngredientStore = create<IngredientState>((set) => ({
  ingredients: [],
  isLoading: false,
  error: null,

  // 食材リストを設定
  setIngredients: (ingredients) => set({ ingredients }),

  // 食材を追加
  addIngredient: (ingredient) =>
    set((state) => ({
      ingredients: [...state.ingredients, ingredient],
    })),

  // 食材を更新
  updateIngredient: (id, ingredient) =>
    set((state) => ({
      ingredients: state.ingredients.map((item) => (item.id === id ? ingredient : item)),
    })),

  // 食材を削除
  removeIngredient: (id) =>
    set((state) => ({
      ingredients: state.ingredients.filter((item) => item.id !== id),
    })),

  // ローディング状態を設定
  setLoading: (isLoading) => set({ isLoading }),

  // エラーを設定
  setError: (error) => set({ error }),
}));
