import type { Metadata } from 'next';
import { Inter, Noto_Sans_JP } from 'next/font/google';
import type { ReactNode } from 'react';
import Footer from '../src/components/Footer';
import Header from '../src/components/Header';
import './globals.css';

// Optimize fonts with next/font
const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-inter',
});

const notoSansJP = Noto_Sans_JP({
  subsets: ['latin'],
  display: 'swap',
  variable: '--font-noto-sans-jp',
  weight: ['400', '500', '700'],
});

export const metadata: Metadata = {
  title: 'DinnerDecider',
  description: 'Decide what to cook based on your ingredients',
};

/**
 * Root layout component for Next.js App Router
 * Provides the HTML structure and global providers
 */
export default function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <html lang="ja" className={`${inter.variable} ${notoSansJP.variable}`}>
      <body className="antialiased font-sans">
        <div className="flex flex-col min-h-screen bg-gray-50">
          <Header />

          <main
            className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8"
            id="main-content"
          >
            {children}
          </main>

          <Footer />
        </div>
      </body>
    </html>
  );
}
