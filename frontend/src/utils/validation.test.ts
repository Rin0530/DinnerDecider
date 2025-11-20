import { describe, expect, it } from 'vitest';
import { validateIngredient } from './validation';

describe('validateIngredient', () => {
  describe('name validation', () => {
    it('should return error when name is empty', () => {
      const result = validateIngredient({ name: '' });
      expect(result.name).toBe('食材名は必須です');
    });

    it('should return error when name is only whitespace', () => {
      const result = validateIngredient({ name: '   ' });
      expect(result.name).toBe('食材名は必須です');
    });

    it('should return error when name exceeds 255 characters', () => {
      const longName = 'a'.repeat(256);
      const result = validateIngredient({ name: longName });
      expect(result.name).toBe('食材名は255文字以内で入力してください');
    });

    it('should not return error for valid name', () => {
      const result = validateIngredient({ name: 'にんじん' });
      expect(result.name).toBeUndefined();
    });
  });

  describe('quantity validation', () => {
    it('should return error when quantity exceeds 100 characters', () => {
      const longQuantity = 'a'.repeat(101);
      const result = validateIngredient({ name: 'test', quantity: longQuantity });
      expect(result.quantity).toBe('数量は100文字以内で入力してください');
    });

    it('should not return error for valid quantity', () => {
      const result = validateIngredient({ name: 'test', quantity: '2本' });
      expect(result.quantity).toBeUndefined();
    });

    it('should not return error when quantity is undefined', () => {
      const result = validateIngredient({ name: 'test' });
      expect(result.quantity).toBeUndefined();
    });
  });

  describe('expiration_date validation', () => {
    it('should return error for invalid date format', () => {
      const result = validateIngredient({ name: 'test', expiration_date: '2025/11/01' });
      expect(result.expiration_date).toBe('日付はYYYY-MM-DD形式で入力してください');
    });

    it('should return error for invalid date', () => {
      const result = validateIngredient({ name: 'test', expiration_date: '2025-13-01' });
      expect(result.expiration_date).toBe('有効な日付を入力してください');
    });

    it('should not return error for valid date', () => {
      const result = validateIngredient({ name: 'test', expiration_date: '2025-11-01' });
      expect(result.expiration_date).toBeUndefined();
    });

    it('should not return error when expiration_date is undefined', () => {
      const result = validateIngredient({ name: 'test' });
      expect(result.expiration_date).toBeUndefined();
    });
  });

  describe('multiple validation errors', () => {
    it('should return multiple errors when multiple fields are invalid', () => {
      const result = validateIngredient({
        name: '',
        quantity: 'a'.repeat(101),
        expiration_date: 'invalid',
      });
      expect(result.name).toBeDefined();
      expect(result.quantity).toBeDefined();
      expect(result.expiration_date).toBeDefined();
    });
  });
});
