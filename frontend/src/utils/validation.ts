import type { IngredientInput } from '../types';

// バリデーションエラーの型定義
export interface ValidationErrors {
  name?: string;
  quantity?: string;
  purchase_date?: string;
}

/**
 * 食材データのバリデーション
 * @param data 検証する食材データ
 * @returns バリデーションエラーのオブジェクト
 */
export const validateIngredient = (data: IngredientInput): ValidationErrors => {
  const errors: ValidationErrors = {};

  // 食材名のバリデーション: 必須、1-255文字
  if (!data.name || data.name.trim() === '') {
    errors.name = '食材名は必須です';
  } else if (data.name.length > 255) {
    errors.name = '食材名は255文字以内で入力してください';
  }

  // 数量のバリデーション: 任意、1-100文字
  if (data.quantity && data.quantity.length > 100) {
    errors.quantity = '数量は100文字以内で入力してください';
  }

  // 購入日のバリデーション: 任意、YYYY-MM-DD形式
  if (data.purchase_date) {
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
    if (!dateRegex.test(data.purchase_date)) {
      errors.purchase_date = '日付はYYYY-MM-DD形式で入力してください';
    } else {
      // 有効な日付かチェック
      const date = new Date(data.purchase_date);
      if (Number.isNaN(date.getTime())) {
        errors.purchase_date = '有効な日付を入力してください';
      }
    }
  }

  return errors;
};
