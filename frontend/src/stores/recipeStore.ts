import { create } from 'zustand';
import type { Recipe } from '../types';

// LocalStorageのキー
const FAVORITES_KEY = 'dinner-decider-favorites';

/**
 * お気に入りをLocalStorageに保存
 */
const saveFavoritesToStorage = (favorites: Recipe[]): void => {
  // Check if running in browser (client-side)
  if (typeof window === 'undefined') return;
  
  try {
    localStorage.setItem(FAVORITES_KEY, JSON.stringify(favorites));
  } catch (error) {
    console.error('Failed to save favorites to localStorage:', error);
  }
};

/**
 * LocalStorageからお気に入りを読み込み
 */
const loadFavoritesFromStorage = (): Recipe[] => {
  // Check if running in browser (client-side)
  if (typeof window === 'undefined') return [];
  
  try {
    const stored = localStorage.getItem(FAVORITES_KEY);
    return stored ? JSON.parse(stored) : [];
  } catch (error) {
    console.error('Failed to load favorites from localStorage:', error);
    return [];
  }
};

// レシピストアの状態とアクションの型定義
interface RecipeState {
  recipes: Recipe[]; // 提案されたレシピリスト
  favorites: Recipe[]; // お気に入りレシピリスト
  isLoading: boolean; // ローディング状態
  error: string | null; // エラーメッセージ
  setRecipes: (recipes: Recipe[]) => void; // レシピリストを設定
  addFavorite: (recipe: Recipe) => void; // お気に入りに追加
  removeFavorite: (recipeName: string) => void; // お気に入りから削除
  setLoading: (isLoading: boolean) => void; // ローディング状態を設定
  setError: (error: string | null) => void; // エラーを設定
}

/**
 * レシピストア
 * Zustandを使用したレシピとお気に入りの状態管理
 */
export const useRecipeStore = create<RecipeState>((set) => ({
  recipes: [],
  favorites: loadFavoritesFromStorage(),
  isLoading: false,
  error: null,

  // レシピリストを設定
  setRecipes: (recipes) => set({ recipes }),

  // お気に入りに追加
  addFavorite: (recipe) =>
    set((state) => {
      // 既にお気に入りに存在するかチェック
      const exists = state.favorites.some((fav) => fav.name === recipe.name);
      if (exists) {
        return state;
      }

      const newFavorites = [...state.favorites, recipe];
      saveFavoritesToStorage(newFavorites);
      return { favorites: newFavorites };
    }),

  // お気に入りから削除
  removeFavorite: (recipeName) =>
    set((state) => {
      const newFavorites = state.favorites.filter((fav) => fav.name !== recipeName);
      saveFavoritesToStorage(newFavorites);
      return { favorites: newFavorites };
    }),

  // ローディング状態を設定
  setLoading: (isLoading) => set({ isLoading }),

  // エラーを設定
  setError: (error) => set({ error }),
}));
