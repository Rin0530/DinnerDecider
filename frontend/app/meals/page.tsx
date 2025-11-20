'use client';

import MealCard from '../../src/components/MealCard';
import { useRecipeStore } from '../../src/stores/recipeStore';

/**
 * お気に入り献立ページコンポーネント
 * 保存したお気に入りレシピを一覧表示
 */
export default function MealsPage() {
  const { favorites, removeFavorite } = useRecipeStore();

  // お気に入り削除時のハンドラー
  const handleRemoveFavorite = (recipeName: string) => {
    if (window.confirm(`「${recipeName}」をお気に入りから削除しますか？`)) {
      removeFavorite(recipeName);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">お気に入り献立</h1>
        <p className="text-gray-600">保存した献立の一覧です</p>
      </div>

      {/* Empty State */}
      {favorites.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg mb-4">お気に入りの献立がまだありません</p>
          <p className="text-gray-400 text-sm">
            献立提案ページで気に入った献立を保存してみましょう
          </p>
        </div>
      ) : (
        <div className="space-y-6">
          {favorites.map((recipe, index) => (
            <div key={`${recipe.name}-${index}`} className="relative">
              <MealCard
                recipe={recipe}
                onSaveFavorite={() => {
                  // Already in favorites, no action needed
                }}
                isFavorite={true}
              />

              {/* Remove from Favorites Button */}
              <div className="mt-4">
                <button
                  type="button"
                  onClick={() => handleRemoveFavorite(recipe.name)}
                  className="w-full sm:w-auto px-6 py-3 text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 font-medium transition-colors min-h-[44px]"
                  aria-label={`${recipe.name}をお気に入りから削除`}
                >
                  お気に入りから削除
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
