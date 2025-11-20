import { useCallback } from 'react';
import { ingredientApi } from '../services/ingredientApi';
import { useIngredientStore } from '../stores/ingredientStore';
import type { ApiError, Ingredient, IngredientInput } from '../types';

// useIngredientsフックの戻り値の型定義
export interface UseIngredientsReturn {
  ingredients: Ingredient[]; // 食材リスト
  isLoading: boolean; // ローディング状態
  error: string | null; // エラーメッセージ
  fetchIngredients: () => Promise<void>; // 食材を取得
  addIngredient: (data: IngredientInput) => Promise<void>; // 食材を追加
  updateIngredient: (id: number, data: IngredientInput) => Promise<void>; // 食材を更新
  deleteIngredient: (id: number) => Promise<void>; // 食材を削除
}

/**
 * 食材管理カスタムフック
 * 食材の取得、追加、更新、削除を行う
 */
export const useIngredients = (): UseIngredientsReturn => {
  const {
    ingredients,
    isLoading,
    error,
    setIngredients,
    addIngredient: addToStore,
    updateIngredient: updateInStore,
    removeIngredient,
    setLoading,
    setError,
  } = useIngredientStore();

  // 食材を取得
  const fetchIngredients = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await ingredientApi.getAll();
      setIngredients(data);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.message || '食材の取得に失敗しました。');
    } finally {
      setLoading(false);
    }
  }, [setIngredients, setLoading, setError]);

  // 食材を追加
  const addIngredient = useCallback(
    async (data: IngredientInput) => {
      setLoading(true);
      setError(null);

      try {
        const newIngredient = await ingredientApi.create(data);
        addToStore(newIngredient);
      } catch (err) {
        const apiError = err as ApiError;
        setError(apiError.message || '食材の追加に失敗しました。');
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [addToStore, setLoading, setError],
  );

  // 食材を更新
  const updateIngredient = useCallback(
    async (id: number, data: IngredientInput) => {
      setLoading(true);
      setError(null);

      try {
        const updatedIngredient = await ingredientApi.update(id, data);
        updateInStore(id, updatedIngredient);
      } catch (err) {
        const apiError = err as ApiError;
        setError(apiError.message || '食材の更新に失敗しました。');
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [updateInStore, setLoading, setError],
  );

  // 食材を削除
  const deleteIngredient = useCallback(
    async (id: number) => {
      setLoading(true);
      setError(null);

      try {
        await ingredientApi.delete(id);
        removeIngredient(id);
      } catch (err) {
        const apiError = err as ApiError;
        setError(apiError.message || '食材の削除に失敗しました。');
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [removeIngredient, setLoading, setError],
  );

  return {
    ingredients,
    isLoading,
    error,
    fetchIngredients,
    addIngredient,
    updateIngredient,
    deleteIngredient,
  };
};
