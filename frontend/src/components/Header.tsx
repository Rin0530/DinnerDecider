'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

/**
 * ヘッダーコンポーネント
 * アプリケーションのナビゲーションを提供
 */
const Header = () => {
  const pathname = usePathname();

  // ナビゲーションアイテムの定義
  const navItems = [
    { path: '/', label: 'ホーム' },
    { path: '/meals', label: 'お気に入り' },
    { path: '/settings', label: '設定' },
  ];

  // 現在のパスがアクティブかどうかを判定
  const isActive = (path: string) => pathname === path;

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <Link
            href="/"
            className="text-xl sm:text-2xl font-bold text-blue-600 hover:text-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded"
            aria-label="DinnerDecider ホームページへ"
          >
            DinnerDecider
          </Link>

          <nav aria-label="メインナビゲーション">
            <ul className="flex space-x-2 sm:space-x-4">
              {navItems.map((item) => (
                <li key={item.path}>
                  <Link
                    href={item.path}
                    className={`
                      px-3 py-2 rounded-md text-sm font-medium transition-colors
                      focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                      ${
                        isActive(item.path)
                          ? 'bg-blue-100 text-blue-700'
                          : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                      }
                    `}
                    aria-current={isActive(item.path) ? 'page' : undefined}
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>
    </header>
  );
};

export default Header;
