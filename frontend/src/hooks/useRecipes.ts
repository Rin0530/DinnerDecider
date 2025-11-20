import { useCallback } from 'react';
import { recipeApi } from '../services/recipeApi';
import { useRecipeStore } from '../stores/recipeStore';
import type { ApiError, Recipe } from '../types';

// useRecipesフックの戻り値の型定義
export interface UseRecipesReturn {
  recipes: Recipe[]; // 提案されたレシピリスト
  favorites: Recipe[]; // お気に入りレシピリスト
  isLoading: boolean; // ローディング状態
  error: string | null; // エラーメッセージ
  requestSuggestion: () => Promise<void>; // 献立提案をリクエスト
  saveFavorite: (recipe: Recipe) => void; // お気に入りに保存
  removeFavorite: (recipeName: string) => void; // お気に入りから削除
  getFavorites: () => Recipe[]; // お気に入りを取得
}

/**
 * レシピ管理カスタムフック
 * レシピの提案、お気に入りの管理を行う
 */
export const useRecipes = (): UseRecipesReturn => {
  const {
    recipes,
    favorites,
    isLoading,
    error,
    setRecipes,
    addFavorite,
    removeFavorite: removeFromFavorites,
    setLoading,
    setError,
  } = useRecipeStore();

  // 献立提案をリクエスト
  const requestSuggestion = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await recipeApi.requestSuggestion();
      setRecipes(response.suggestions);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || '献立の提案に失敗しました。');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [setRecipes, setLoading, setError]);

  // お気に入りに保存
  const saveFavorite = useCallback(
    (recipe: Recipe) => {
      addFavorite(recipe);
    },
    [addFavorite],
  );

  // お気に入りから削除
  const removeFavorite = useCallback(
    (recipeName: string) => {
      removeFromFavorites(recipeName);
    },
    [removeFromFavorites],
  );

  // お気に入りを取得
  const getFavorites = useCallback(() => {
    return favorites;
  }, [favorites]);

  return {
    recipes,
    favorites,
    isLoading,
    error,
    requestSuggestion,
    saveFavorite,
    removeFavorite,
    getFavorites,
  };
};
