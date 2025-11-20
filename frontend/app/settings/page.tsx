/**
 * 設定ページコンポーネント
 * アプリケーションの設定を管理（今後実装予定）
 */
export default function SettingsPage() {
  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">設定</h1>
        <p className="text-gray-600">アプリケーションの設定を管理します</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <p className="text-gray-500 text-center py-8">設定項目は今後追加予定です</p>
      </div>
    </div>
  );
}
