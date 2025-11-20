# DinnerDecider Frontend

冷蔵庫の食材から献立を提案する Web アプリケーションのフロントエンド

## 技術スタック

- React 18 + TypeScript
- Vite (ビルドツール)
- React Router (ルーティング)
- Zustand (状態管理)
- Tailwind CSS (スタイリング)
- Vitest + React Testing Library (テスト)
- Biome (Lint/Format)

## セットアップ

### 前提条件

- Node.js 18 以上
- npm または yarn

### インストール

```bash
# 依存関係のインストール
pnpm install

# 環境変数の設定
cp .env.example .env
# .env ファイルを編集して API_BASE_URL を設定
```

### 開発サーバーの起動

```bash
npm run dev
```

ブラウザで http://localhost:5173 を開きます。

## スクリプト

- `npm run dev` - 開発サーバーを起動
- `npm run build` - プロダクションビルド
- `npm run preview` - ビルドしたアプリをプレビュー
- `npm run lint` - Biome でコードをチェック
- `npm run format` - Biome でコードをフォーマット
- `npm run type-check` - TypeScript の型チェック
- `npm run test` - テストを実行
- `npm run test:watch` - テストをウォッチモードで実行

## ディレクトリ構造

```
src/
├── main.tsx                 # アプリケーションエントリポイント
├── App.tsx                  # ルートコンポーネント
├── pages/                   # ページコンポーネント
├── components/              # 再利用可能なコンポーネント
├── layouts/                 # レイアウトコンポーネント
├── hooks/                   # カスタムフック
├── stores/                  # Zustand ストア
├── services/                # API クライアント
├── types/                   # TypeScript 型定義
├── utils/                   # ユーティリティ関数
├── styles/                  # グローバルスタイル
└── assets/                  # 静的アセット
```

## 環境変数

`.env` ファイルで以下の環境変数を設定します：

- `VITE_API_BASE_URL` - バックエンド API のベース URL (デフォルト:
  http://localhost:8080)

## コーディング規約

- Biome を使用してコードの品質を保ちます
- TypeScript の strict モードを有効にしています
- アクセシビリティ (a11y) を重視した実装を心がけます
- モバイルファーストのレスポンシブデザインを採用します

## ライセンス

Private
