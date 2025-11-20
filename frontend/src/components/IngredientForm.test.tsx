import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@/src/test/utils';
import userEvent from '@testing-library/user-event';
import { IngredientForm } from './IngredientForm';
import type { Ingredient } from '@/src/types';

describe('IngredientForm', () => {
  const mockOnSubmit = vi.fn();
  const mockOnCancel = vi.fn();

  beforeEach(() => {
    mockOnSubmit.mockClear();
    mockOnCancel.mockClear();
  });

  describe('rendering', () => {
    it('should render add form when no ingredient is provided', () => {
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);
      expect(screen.getByText('食材を追加')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '追加' })).toBeInTheDocument();
    });

    it('should render edit form when ingredient is provided', () => {
      const ingredient: Ingredient = {
        id: 1,
        name: 'にんじん',
        quantity: '2本',
        purchase_date: '2025-11-01',
      };
      render(
        <IngredientForm ingredient={ingredient} onSubmit={mockOnSubmit} onCancel={mockOnCancel} />,
      );
      expect(screen.getByText('食材を編集')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '更新' })).toBeInTheDocument();
    });

    it('should pre-fill form with ingredient data', () => {
      const ingredient: Ingredient = {
        id: 1,
        name: 'にんじん',
        quantity: '2本',
        purchase_date: '2025-11-01',
      };
      render(
        <IngredientForm ingredient={ingredient} onSubmit={mockOnSubmit} onCancel={mockOnCancel} />,
      );

      expect(screen.getByLabelText(/食材名/)).toHaveValue('にんじん');
      expect(screen.getByLabelText(/数量/)).toHaveValue('2本');
      expect(screen.getByLabelText(/購入日/)).toHaveValue('2025-11-01');
    });
  });

  describe('form submission', () => {
    it('should call onSubmit with form data when valid', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.type(screen.getByLabelText(/食材名/), 'にんじん');
      await user.type(screen.getByLabelText(/数量/), '2本');
      await user.type(screen.getByLabelText(/購入日/), '2025-11-01');

      await user.click(screen.getByRole('button', { name: '追加' }));

      expect(mockOnSubmit).toHaveBeenCalledWith({
        name: 'にんじん',
        quantity: '2本',
        purchase_date: '2025-11-01',
      });
    });

    it('should not call onSubmit when name is empty', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.click(screen.getByRole('button', { name: '追加' }));

      expect(mockOnSubmit).not.toHaveBeenCalled();
      expect(screen.getByText('食材名は必須です')).toBeInTheDocument();
    });

    it('should submit with only required fields', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.type(screen.getByLabelText(/食材名/), 'にんじん');
      await user.click(screen.getByRole('button', { name: '追加' }));

      expect(mockOnSubmit).toHaveBeenCalledWith({
        name: 'にんじん',
        quantity: '',
        purchase_date: '',
      });
    });
  });

  describe('validation', () => {
    it('should show error when name exceeds 255 characters', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      const longName = 'a'.repeat(256);
      await user.type(screen.getByLabelText(/食材名/), longName);
      await user.click(screen.getByRole('button', { name: '追加' }));

      expect(screen.getByText('食材名は255文字以内で入力してください')).toBeInTheDocument();
      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it('should show error when quantity exceeds 100 characters', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.type(screen.getByLabelText(/食材名/), 'にんじん');
      const longQuantity = 'a'.repeat(101);
      await user.type(screen.getByLabelText(/数量/), longQuantity);
      await user.click(screen.getByRole('button', { name: '追加' }));

      expect(screen.getByText('数量は100文字以内で入力してください')).toBeInTheDocument();
      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it('should clear error when user starts typing', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      // Trigger validation error
      await user.click(screen.getByRole('button', { name: '追加' }));
      expect(screen.getByText('食材名は必須です')).toBeInTheDocument();

      // Start typing
      await user.type(screen.getByLabelText(/食材名/), 'に');

      // Error should be cleared
      await waitFor(() => {
        expect(screen.queryByText('食材名は必須です')).not.toBeInTheDocument();
      });
    });
  });

  describe('keyboard navigation', () => {
    it('should call onCancel when Escape key is pressed', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      const nameInput = screen.getByLabelText(/食材名/);
      await user.click(nameInput);
      await user.keyboard('{Escape}');

      expect(mockOnCancel).toHaveBeenCalled();
    });
  });

  describe('cancel button', () => {
    it('should call onCancel when cancel button is clicked', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.click(screen.getByRole('button', { name: 'キャンセル' }));

      expect(mockOnCancel).toHaveBeenCalled();
    });
  });

  describe('accessibility', () => {
    it('should have proper aria attributes', () => {
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      const nameInput = screen.getByLabelText(/食材名/);
      expect(nameInput).toHaveAttribute('aria-required', 'true');
    });

    it('should have aria-invalid when field has error', async () => {
      const user = userEvent.setup();
      render(<IngredientForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

      await user.click(screen.getByRole('button', { name: '追加' }));

      const nameInput = screen.getByLabelText(/食材名/);
      expect(nameInput).toHaveAttribute('aria-invalid', 'true');
    });
  });
});
