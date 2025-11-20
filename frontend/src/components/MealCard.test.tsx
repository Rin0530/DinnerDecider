import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@/src/test/utils';
import userEvent from '@testing-library/user-event';
import MealCard from './MealCard';
import type { Recipe } from '@/src/types';

describe('MealCard', () => {
  const mockRecipe: Recipe = {
    name: 'カレーライス',
    steps: ['野菜を切る', 'カレールーを入れる', '煮込む'],
    missing_items: ['カレールー', '肉'],
  };

  it('should render recipe name', () => {
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} />);
    expect(screen.getByText('カレーライス')).toBeInTheDocument();
  });

  it('should render cooking steps', () => {
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} />);
    expect(screen.getByText('野菜を切る')).toBeInTheDocument();
    expect(screen.getByText('カレールーを入れる')).toBeInTheDocument();
    expect(screen.getByText('煮込む')).toBeInTheDocument();
  });

  it('should render missing ingredients', () => {
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} />);
    expect(screen.getByText('カレールー')).toBeInTheDocument();
    expect(screen.getByText('肉')).toBeInTheDocument();
  });

  it('should call onSaveFavorite when save button is clicked', async () => {
    const user = userEvent.setup();
    const onSaveFavorite = vi.fn();
    render(<MealCard recipe={mockRecipe} onSaveFavorite={onSaveFavorite} />);

    const saveButton = screen.getByRole('button', { name: /お気に入りに保存/ });
    await user.click(saveButton);

    expect(onSaveFavorite).toHaveBeenCalledWith(mockRecipe);
  });

  it('should disable button after saving', async () => {
    const user = userEvent.setup();
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} />);

    const saveButton = screen.getByRole('button', { name: /お気に入りに保存/ });
    await user.click(saveButton);

    expect(saveButton).toBeDisabled();
    expect(screen.getByText('保存済み')).toBeInTheDocument();
  });

  it('should show confirmation message after saving', async () => {
    const user = userEvent.setup();
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} />);

    const saveButton = screen.getByRole('button', { name: /お気に入りに保存/ });
    await user.click(saveButton);

    expect(screen.getByText('お気に入りに保存しました')).toBeInTheDocument();
  });

  it('should render as already saved when isFavorite is true', () => {
    render(<MealCard recipe={mockRecipe} onSaveFavorite={() => {}} isFavorite={true} />);

    const saveButton = screen.getByRole('button');
    expect(saveButton).toBeDisabled();
    expect(screen.getByText('保存済み')).toBeInTheDocument();
  });

  it('should not render missing items section when there are no missing items', () => {
    const recipeWithoutMissing: Recipe = {
      name: 'サラダ',
      steps: ['野菜を洗う', '切る'],
      missing_items: [],
    };
    render(<MealCard recipe={recipeWithoutMissing} onSaveFavorite={() => {}} />);

    expect(screen.queryByText('不足している食材')).not.toBeInTheDocument();
  });
});
