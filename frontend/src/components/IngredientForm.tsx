'use client';

import { type FormEvent, useEffect, useState } from 'react';
import type { Ingredient, IngredientInput } from '../types';
import { type ValidationErrors, validateIngredient } from '../utils/validation';

// 食材フォームのプロパティ
interface IngredientFormProps {
  ingredient?: Ingredient; // 編集時の食材データ
  onSubmit: (data: IngredientInput) => void; // 送信時のコールバック
  onCancel: () => void; // キャンセル時のコールバック
}

/**
 * 食材の追加・編集フォームコンポーネント
 */
export const IngredientForm = ({ ingredient, onSubmit, onCancel }: IngredientFormProps) => {
  // フォームデータの状態管理
  const [formData, setFormData] = useState<IngredientInput>(() => {
    // 日付をYYYY-MM-DD形式に変換（ISO形式の場合は日付部分のみ抽出）
    let purchaseDate = ingredient?.purchase_date || '';
    if (purchaseDate && purchaseDate.includes('T')) {
      purchaseDate = purchaseDate.split('T')[0];
    }

    return {
      name: ingredient?.name || '',
      quantity: ingredient?.quantity || '',
      purchase_date: purchaseDate,
    };
  });

  // バリデーションエラーの状態管理
  const [errors, setErrors] = useState<ValidationErrors>({});
  // フィールドがタッチされたかの状態管理
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  // 編集時に食材データをフォームに反映
  useEffect(() => {
    if (ingredient) {
      // 日付をYYYY-MM-DD形式に変換（ISO形式の場合は日付部分のみ抽出）
      let purchaseDate = ingredient.purchase_date || '';
      if (purchaseDate && purchaseDate.includes('T')) {
        purchaseDate = purchaseDate.split('T')[0];
      }

      setFormData({
        name: ingredient.name,
        quantity: ingredient.quantity || '',
        purchase_date: purchaseDate,
      });
    }
  }, [ingredient]);

  // フィールド値変更時のハンドラー
  const handleChange = (field: keyof IngredientInput, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));

    // ユーザーが入力を開始したらエラーをクリア
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  // フィールドからフォーカスが外れた時のハンドラー
  const handleBlur = (field: keyof IngredientInput) => {
    setTouched((prev) => ({ ...prev, [field]: true }));

    // フォーカスが外れた時にバリデーション実行
    const validationErrors = validateIngredient(formData);
    if (validationErrors[field]) {
      setErrors((prev) => ({ ...prev, [field]: validationErrors[field] }));
    }
  };

  // フォーム送信時のハンドラー
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();

    // 全フィールドをタッチ済みにマーク
    setTouched({ name: true, quantity: true, purchase_date: true });

    // 全フィールドをバリデーション
    const validationErrors = validateIngredient(formData);
    setErrors(validationErrors);

    // エラーがなければ送信
    if (Object.keys(validationErrors).length === 0) {
      onSubmit(formData);
    }
  };

  // Escapeキーでキャンセル
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      onCancel();
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      onKeyDown={handleKeyDown}
      className="bg-white p-6 rounded-lg shadow-md"
      aria-label={ingredient ? '食材編集フォーム' : '食材追加フォーム'}
    >
      <h2 className="text-xl font-bold mb-4">{ingredient ? '食材を編集' : '食材を追加'}</h2>

      <div className="mb-4">
        <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
          食材名 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          id="name"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          onBlur={() => handleBlur('name')}
          className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${touched.name && errors.name ? 'border-red-500' : 'border-gray-300'
            }`}
          aria-required="true"
          aria-invalid={touched.name && !!errors.name}
          aria-describedby={touched.name && errors.name ? 'name-error' : undefined}
        />
        {touched.name && errors.name && (
          <p id="name-error" className="text-red-500 text-sm mt-1" role="alert">
            {errors.name}
          </p>
        )}
      </div>

      <div className="mb-4">
        <label htmlFor="quantity" className="block text-sm font-medium text-gray-700 mb-1">
          数量
        </label>
        <input
          type="text"
          id="quantity"
          value={formData.quantity}
          onChange={(e) => handleChange('quantity', e.target.value)}
          onBlur={() => handleBlur('quantity')}
          className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${touched.quantity && errors.quantity ? 'border-red-500' : 'border-gray-300'
            }`}
          placeholder="例: 2個"
          aria-invalid={touched.quantity && !!errors.quantity}
          aria-describedby={touched.quantity && errors.quantity ? 'quantity-error' : undefined}
        />
        {touched.quantity && errors.quantity && (
          <p id="quantity-error" className="text-red-500 text-sm mt-1" role="alert">
            {errors.quantity}
          </p>
        )}
      </div>

      <div className="mb-6">
        <label htmlFor="purchase_date" className="block text-sm font-medium text-gray-700 mb-1">
          購入日
        </label>
        <input
          type="date"
          id="purchase_date"
          value={formData.purchase_date}
          onChange={(e) => handleChange('purchase_date', e.target.value)}
          onBlur={() => handleBlur('purchase_date')}
          className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${touched.purchase_date && errors.purchase_date ? 'border-red-500' : 'border-gray-300'
            }`}
          aria-invalid={touched.purchase_date && !!errors.purchase_date}
          aria-describedby={
            touched.purchase_date && errors.purchase_date ? 'purchase-error' : undefined
          }
        />
        {touched.purchase_date && errors.purchase_date && (
          <p id="purchase-error" className="text-red-500 text-sm mt-1" role="alert">
            {errors.purchase_date}
          </p>
        )}
      </div>

      <div className="flex gap-3 justify-end">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 min-w-[88px] min-h-[44px]"
        >
          キャンセル
        </button>
        <button
          type="submit"
          className="px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[88px] min-h-[44px]"
        >
          {ingredient ? '更新' : '追加'}
        </button>
      </div>
    </form>
  );
};
