'use client';

import { useRouter } from 'next/navigation';
import MealCard from '../../src/components/MealCard';
import { useRecipeStore } from '../../src/stores/recipeStore';
import type { Recipe } from '../../src/types';

/**
 * 献立提案結果ページコンポーネント
 * AIが提案した献立を表示
 */
export default function ResultPage() {
  const router = useRouter();
  const { recipes, addFavorite, favorites } = useRecipeStore();

  // お気に入り保存時のハンドラー
  const handleSaveFavorite = (recipe: Recipe) => {
    addFavorite(recipe);
  };

  // ホームに戻るハンドラー
  const handleBackToHome = () => {
    router.push('/');
  };

  // レシピが既にお気に入りに登録されているかチェック
  const isFavorite = (recipeName: string): boolean => {
    return favorites.some((fav) => fav.name === recipeName);
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">献立提案</h1>
        <p className="text-gray-600">あなたの食材から提案された献立です</p>
      </div>

      {/* Back to Home Button */}
      <div className="mb-6">
        <button
          type="button"
          onClick={handleBackToHome}
          className="px-4 py-2 text-blue-600 bg-white border border-blue-600 rounded-lg hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 min-h-[44px]"
          aria-label="ホームに戻る"
        >
          ← ホームに戻る
        </button>
      </div>

      {/* No Recipes Message */}
      {recipes.length === 0 && (
        <output
          className="p-8 bg-gray-50 border border-gray-200 rounded-lg text-center block"
          aria-live="polite"
        >
          <p className="text-gray-700 text-lg mb-2">献立の提案がありません</p>
          <p className="text-gray-600 text-sm mb-4">食材を登録して、献立を提案してもらいましょう</p>
          <button
            type="button"
            onClick={handleBackToHome}
            className="px-6 py-3 text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 min-h-[44px]"
            aria-label="ホームに戻る"
          >
            ホームに戻る
          </button>
        </output>
      )}

      {/* Recipe List */}
      {recipes.length > 0 && (
        <ul className="grid grid-cols-1 md:grid-cols-2 gap-6 list-none" aria-label="献立リスト">
          {recipes.map((recipe, index) => (
            <li key={`${recipe.name}-${index}`}>
              <MealCard
                recipe={recipe}
                onSaveFavorite={handleSaveFavorite}
                isFavorite={isFavorite(recipe.name)}
              />
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
