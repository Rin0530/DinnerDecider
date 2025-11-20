import { describe, expect, it } from 'vitest';
import { render, screen } from '@/src/test/utils';
import LoadingSpinner from './LoadingSpinner';

describe('LoadingSpinner', () => {
  it('should render with default size', () => {
    render(<LoadingSpinner />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('should render with custom message', () => {
    render(<LoadingSpinner message="データを読み込んでいます..." />);
    expect(screen.getByText('データを読み込んでいます...')).toBeInTheDocument();
  });

  it('should have aria-live and aria-busy attributes', () => {
    const { container } = render(<LoadingSpinner />);
    const output = container.querySelector('output');
    expect(output).toHaveAttribute('aria-live', 'polite');
    expect(output).toHaveAttribute('aria-busy', 'true');
  });

  it('should render with small size', () => {
    const { container } = render(<LoadingSpinner size="sm" />);
    const spinner = container.querySelector('.h-6.w-6');
    expect(spinner).toBeInTheDocument();
  });

  it('should render with large size', () => {
    const { container } = render(<LoadingSpinner size="lg" />);
    const spinner = container.querySelector('.h-16.w-16');
    expect(spinner).toBeInTheDocument();
  });
});
