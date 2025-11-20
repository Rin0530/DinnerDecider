import LoadingSpinner from '../src/components/LoadingSpinner';

/**
 * Loading component for Next.js App Router
 * Displayed while a page is loading
 */
export default function Loading() {
  return (
    <div className="flex items-center justify-center min-h-[400px]">
      <LoadingSpinner size="lg" message="読み込み中..." />
    </div>
  );
}
