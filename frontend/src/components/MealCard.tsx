'use client';

import { useState } from 'react';
import type { Recipe } from '../types';

// 献立カードのプロパティ
interface MealCardProps {
  recipe: Recipe; // 表示するレシピ
  onSaveFavorite: (recipe: Recipe) => void; // お気に入り保存時のコールバック
  isFavorite?: boolean; // 既にお気に入りかどうか
}

/**
 * 献立カードコンポーネント
 * レシピの詳細を表示し、お気に入りに保存する機能を提供
 */
const MealCard = ({ recipe, onSaveFavorite, isFavorite = false }: MealCardProps) => {
  // 保存済み状態の管理
  const [isSaved, setIsSaved] = useState(isFavorite);
  // 確認メッセージの表示状態
  const [showConfirmation, setShowConfirmation] = useState(false);

  // お気に入り保存時のハンドラー
  const handleSaveFavorite = () => {
    onSaveFavorite(recipe);
    setIsSaved(true);
    setShowConfirmation(true);

    // 2秒後に確認メッセージを非表示
    setTimeout(() => {
      setShowConfirmation(false);
    }, 2000);
  };

  return (
    <article
      className="bg-white rounded-lg shadow-md p-4 sm:p-6 hover:shadow-lg transition-shadow"
      aria-labelledby={`recipe-${recipe.name.replace(/\s+/g, '-')}`}
    >
      {/* Recipe Name */}
      <h3
        id={`recipe-${recipe.name.replace(/\s+/g, '-')}`}
        className="text-xl sm:text-2xl font-bold text-gray-900 mb-4"
      >
        {recipe.name}
      </h3>

      {/* Cooking Steps */}
      <section className="mb-4">
        <h4 className="text-lg font-semibold text-gray-800 mb-2">調理手順</h4>
        <ol className="list-decimal list-inside space-y-2">
          {recipe.steps.map((step) => (
            <li key={step} className="text-gray-700 leading-relaxed">
              {step}
            </li>
          ))}
        </ol>
      </section>

      {/* Missing Ingredients */}
      {recipe.missing_items.length > 0 && (
        <section className="mb-4">
          <h4 className="text-lg font-semibold text-gray-800 mb-2">不足している食材</h4>
          <ul className="list-disc list-inside space-y-1">
            {recipe.missing_items.map((item) => (
              <li key={item} className="text-red-600">
                {item}
              </li>
            ))}
          </ul>
        </section>
      )}

      {/* Save to Favorites Button */}
      <div className="mt-6">
        <button
          type="button"
          onClick={handleSaveFavorite}
          disabled={isSaved}
          className={`
            w-full sm:w-auto px-6 py-3 rounded-md font-medium transition-colors
            focus:outline-none focus:ring-2 focus:ring-offset-2
            ${
              isSaved
                ? 'bg-gray-300 text-gray-600 cursor-not-allowed'
                : 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500'
            }
          `}
          aria-label={
            isSaved
              ? `${recipe.name}は既にお気に入りに保存されています`
              : `${recipe.name}をお気に入りに保存`
          }
        >
          {isSaved ? '保存済み' : 'お気に入りに保存'}
        </button>

        {/* Confirmation Message */}
        {showConfirmation && (
          <output className="mt-2 text-green-600 font-medium block" aria-live="polite">
            お気に入りに保存しました
          </output>
        )}
      </div>
    </article>
  );
};

export default MealCard;
