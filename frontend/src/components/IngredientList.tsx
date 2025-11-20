'use client';

import { useState } from 'react';
import type { Ingredient } from '../types';

// 食材リストのプロパティ
interface IngredientListProps {
  ingredients: Ingredient[]; // 表示する食材の配列
  onEdit: (ingredient: Ingredient) => void; // 編集時のコールバック
  onDelete: (id: number) => void; // 削除時のコールバック
}

/**
 * 食材リストコンポーネント
 * 登録された食材を一覧表示し、編集・削除機能を提供
 */
export const IngredientList = ({ ingredients, onEdit, onDelete }: IngredientListProps) => {
  // 削除確認ダイアログの表示状態
  const [deleteConfirmId, setDeleteConfirmId] = useState<number | null>(null);

  // 削除ボタンクリック時のハンドラー
  const handleDeleteClick = (id: number) => {
    setDeleteConfirmId(id);
  };

  // 削除確認時のハンドラー
  const handleConfirmDelete = () => {
    if (deleteConfirmId !== null) {
      onDelete(deleteConfirmId);
      setDeleteConfirmId(null);
    }
  };

  // 削除キャンセル時のハンドラー
  const handleCancelDelete = () => {
    setDeleteConfirmId(null);
  };

  // キーボード操作のハンドラー
  const handleKeyDown = (e: React.KeyboardEvent, callback: () => void) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      callback();
    }
  };

  if (ingredients.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <p>食材が登録されていません</p>
        <p className="text-sm mt-2">「食材を追加」ボタンから登録してください</p>
      </div>
    );
  }

  return (
    <>
      <ul className="space-y-3" aria-label="食材一覧">
        {ingredients.map((ingredient) => (
          <li
            key={ingredient.id}
            className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
          >
            <div className="flex justify-between items-start gap-4">
              <div className="flex-1 min-w-0">
                <h3 className="text-lg font-semibold text-gray-900 truncate">{ingredient.name}</h3>
                <div className="mt-1 space-y-1">
                  {ingredient.quantity && (
                    <p className="text-sm text-gray-600">
                      <span className="font-medium">数量:</span> {ingredient.quantity}
                    </p>
                  )}
                  {ingredient.purchase_date && (
                    <p className="text-sm text-gray-600">
                      <span className="font-medium">購入日:</span> {ingredient.purchase_date}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex gap-2 flex-shrink-0">
                <button
                  type="button"
                  onClick={() => onEdit(ingredient)}
                  className="px-3 py-2 text-sm text-blue-600 bg-blue-50 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[60px] min-h-[44px]"
                  aria-label={`${ingredient.name}を編集`}
                >
                  編集
                </button>
                <button
                  type="button"
                  onClick={() => handleDeleteClick(ingredient.id)}
                  className="px-3 py-2 text-sm text-red-600 bg-red-50 rounded-md hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-red-500 min-w-[60px] min-h-[44px]"
                  aria-label={`${ingredient.name}を削除`}
                >
                  削除
                </button>
              </div>
            </div>
          </li>
        ))}
      </ul>

      {/* Delete Confirmation Dialog */}
      {deleteConfirmId !== null && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
          role="dialog"
          aria-modal="true"
          aria-labelledby="delete-dialog-title"
          onClick={handleCancelDelete}
          onKeyDown={(e: React.KeyboardEvent) => {
            if (e.key === 'Escape') handleCancelDelete();
          }}
        >
          <div
            className="bg-white rounded-lg shadow-xl max-w-md w-full p-6"
            onClick={(e: React.MouseEvent) => e.stopPropagation()}
            onKeyDown={(e: React.KeyboardEvent) => e.stopPropagation()}
          >
            <h2 id="delete-dialog-title" className="text-xl font-bold text-gray-900 mb-4">
              削除の確認
            </h2>
            <p className="text-gray-700 mb-6">
              この食材を削除してもよろしいですか？
              <br />
              この操作は取り消せません。
            </p>
            <div className="flex gap-3 justify-end">
              <button
                type="button"
                onClick={handleCancelDelete}
                onKeyDown={(e) => handleKeyDown(e, handleCancelDelete)}
                className="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 min-w-[88px] min-h-[44px]"
              >
                キャンセル
              </button>
              <button
                type="button"
                onClick={handleConfirmDelete}
                onKeyDown={(e) => handleKeyDown(e, handleConfirmDelete)}
                className="px-4 py-2 text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 min-w-[88px] min-h-[44px]"
              >
                削除
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
