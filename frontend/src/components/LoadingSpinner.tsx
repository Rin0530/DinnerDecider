// ローディングスピナーのプロパティ
interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg'; // スピナーのサイズ
  message?: string; // 表示するメッセージ
}

/**
 * ローディングスピナーコンポーネント
 * 処理中の状態を視覚的に表示
 */
const LoadingSpinner = ({ size = 'md', message }: LoadingSpinnerProps) => {
  // サイズに応じたCSSクラス
  const sizeClasses = {
    sm: 'h-6 w-6',
    md: 'h-12 w-12',
    lg: 'h-16 w-16',
  };

  return (
    <output
      className="flex flex-col items-center justify-center"
      aria-live="polite"
      aria-busy="true"
    >
      <div
        className={`${sizeClasses[size]} animate-spin rounded-full border-4 border-gray-200 border-t-blue-600`}
        aria-hidden="true"
      />
      {message && (
        <p className="mt-4 text-sm text-gray-600" aria-label={message}>
          {message}
        </p>
      )}
      <span className="sr-only">読み込み中...</span>
    </output>
  );
};

export default LoadingSpinner;
