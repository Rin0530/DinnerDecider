import { Component, type ReactNode } from 'react';
import { ERROR_MESSAGES } from '../utils/errorMessages';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

/**
 * Error Boundary コンポーネント
 * 予期しないエラーをキャッチしてフォールバック UI を表示する
 */
class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    // エラーが発生したら state を更新
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    // エラーログをコンソールに出力
    console.error('ErrorBoundary caught an error:', error, errorInfo);
  }

  handleReset = (): void => {
    // エラー状態をリセット
    this.setState({
      hasError: false,
      error: null,
    });
  };

  handleReload = (): void => {
    // ページをリロード
    window.location.reload();
  };

  render(): ReactNode {
    if (this.state.hasError) {
      // カスタムフォールバック UI が提供されている場合はそれを使用
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // デフォルトのエラー UI を表示
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
          <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
            <div className="text-center">
              {/* エラーアイコン */}
              <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-red-100 mb-4">
                <svg
                  className="h-8 w-8 text-red-600"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
              </div>

              {/* エラーメッセージ */}
              <h2 className="text-2xl font-bold text-gray-900 mb-2">エラーが発生しました</h2>
              <p className="text-gray-600 mb-6">{ERROR_MESSAGES.UNKNOWN_ERROR}</p>

              {/* エラー詳細（開発環境のみ） */}
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <div className="mb-6 p-4 bg-gray-100 rounded-md text-left">
                  <p className="text-sm font-mono text-gray-800 break-words">
                    {this.state.error.message}
                  </p>
                </div>
              )}

              {/* アクションボタン */}
              <div className="flex flex-col sm:flex-row gap-3 justify-center">
                <button
                  type="button"
                  onClick={this.handleReset}
                  className="px-6 py-3 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
                >
                  もう一度試す
                </button>
                <button
                  type="button"
                  onClick={this.handleReload}
                  className="px-6 py-3 bg-gray-200 text-gray-800 font-medium rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors"
                >
                  ページを再読み込み
                </button>
              </div>

              {/* ホームに戻るリンク */}
              <div className="mt-6">
                <a href="/" className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                  ホームに戻る
                </a>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
