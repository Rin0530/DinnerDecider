'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { IngredientForm } from '../src/components/IngredientForm';
import { IngredientList } from '../src/components/IngredientList';
import LoadingSpinner from '../src/components/LoadingSpinner';
import { useIngredients } from '../src/hooks/useIngredients';
import { useRecipes } from '../src/hooks/useRecipes';
import type { Ingredient, IngredientInput } from '../src/types';

/**
 * ホームページコンポーネント
 * 食材の管理と献立提案のリクエストを行う
 */
export default function HomePage() {
  const router = useRouter();
  
  // 食材管理フック
  const {
    ingredients,
    isLoading: isLoadingIngredients,
    error: ingredientError,
    fetchIngredients,
    addIngredient,
    updateIngredient,
    deleteIngredient,
  } = useIngredients();

  // レシピ管理フック
  const { isLoading: isLoadingRecipes, error: recipeError, requestSuggestion } = useRecipes();

  // フォーム表示状態
  const [showForm, setShowForm] = useState(false);
  // 編集中の食材
  const [editingIngredient, setEditingIngredient] = useState<Ingredient | undefined>(undefined);

  // コンポーネントマウント時に食材を取得
  useEffect(() => {
    fetchIngredients();
  }, [fetchIngredients]);

  // 追加ボタンクリック時のハンドラー
  const handleAddClick = () => {
    setEditingIngredient(undefined);
    setShowForm(true);
  };

  // 編集ボタンクリック時のハンドラー
  const handleEditClick = (ingredient: Ingredient) => {
    setEditingIngredient(ingredient);
    setShowForm(true);
  };

  // フォーム送信時のハンドラー
  const handleFormSubmit = async (data: IngredientInput) => {
    try {
      if (editingIngredient) {
        await updateIngredient(editingIngredient.id, data);
      } else {
        await addIngredient(data);
      }
      setShowForm(false);
      setEditingIngredient(undefined);
    } catch (error) {
      // エラーはフック内で処理済み
      console.error('Failed to save ingredient:', error);
    }
  };

  // フォームキャンセル時のハンドラー
  const handleFormCancel = () => {
    setShowForm(false);
    setEditingIngredient(undefined);
  };

  // 削除ボタンクリック時のハンドラー
  const handleDeleteClick = async (id: number) => {
    try {
      await deleteIngredient(id);
    } catch (error) {
      // エラーはフック内で処理済み
      console.error('Failed to delete ingredient:', error);
    }
  };

  // 献立提案リクエスト時のハンドラー
  const handleRequestSuggestion = async () => {
    if (ingredients.length === 0) {
      return;
    }

    try {
      await requestSuggestion();
      // 提案成功後に結果ページへ遷移
      router.push('/result');
    } catch (error) {
      // エラーはフック内で処理済み
      console.error('Failed to request suggestion:', error);
    }
  };

  const isLoading = isLoadingIngredients || isLoadingRecipes;
  const error = ingredientError || recipeError;

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">冷蔵庫の食材</h1>
        <p className="text-gray-600">食材を登録して、献立を提案してもらいましょう</p>
      </div>

      {/* Error Message */}
      {error && (
        <div
          className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg"
          role="alert"
          aria-live="assertive"
        >
          <p className="text-red-800 font-medium">エラーが発生しました</p>
          <p className="text-red-700 text-sm mt-1">{error}</p>
        </div>
      )}

      {/* Loading State */}
      {isLoading && (
        <div className="mb-6">
          <LoadingSpinner message="処理中..." />
        </div>
      )}

      {/* Add Ingredient Button */}
      {!showForm && (
        <div className="mb-6">
          <button
            type="button"
            onClick={handleAddClick}
            disabled={isLoading}
            className="px-6 py-3 text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed min-h-[44px]"
            aria-label="食材を追加"
          >
            食材を追加
          </button>
        </div>
      )}

      {/* Ingredient Form */}
      {showForm && (
        <div className="mb-6">
          <IngredientForm
            ingredient={editingIngredient}
            onSubmit={handleFormSubmit}
            onCancel={handleFormCancel}
          />
        </div>
      )}

      {/* Ingredient List */}
      <div className="mb-6">
        <IngredientList
          ingredients={ingredients}
          onEdit={handleEditClick}
          onDelete={handleDeleteClick}
        />
      </div>

      {/* Suggest Recipes Button */}
      {ingredients.length > 0 && !showForm && (
        <div className="mt-8 pt-6 border-t border-gray-200">
          <button
            type="button"
            onClick={handleRequestSuggestion}
            disabled={isLoading}
            className="w-full px-6 py-4 text-white bg-green-600 rounded-lg hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-400 disabled:cursor-not-allowed min-h-[44px] text-lg font-semibold"
            aria-label="献立を提案"
          >
            献立を提案
          </button>
          <p className="text-center text-sm text-gray-500 mt-2">
            登録された食材から献立を提案します
          </p>
        </div>
      )}
    </div>
  );
}
