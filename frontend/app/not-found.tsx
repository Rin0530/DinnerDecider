import Link from 'next/link';

/**
 * 404 Not Found page component
 * Displayed when a route is not found
 */
export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] px-4">
      <div className="max-w-md w-full text-center">
        <div className="mb-8">
          <h1 className="text-6xl font-bold text-gray-900 mb-2">404</h1>
          <h2 className="text-2xl font-semibold text-gray-700 mb-4">
            ページが見つかりません
          </h2>
          <p className="text-gray-600">
            お探しのページは存在しないか、移動した可能性があります。
          </p>
        </div>

        <Link
          href="/"
          className="inline-block px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
        >
          ホームに戻る
        </Link>
      </div>
    </div>
  );
}
