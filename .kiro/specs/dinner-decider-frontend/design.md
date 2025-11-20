# Design Document

## Overview

DinnerDecider フロントエンドは、React + TypeScript で構築されるシングルページアプリケーション (SPA) です。Vite をビルドツールとして使用し、モバイルファーストのレスポンシブデザインを実装します。アプリケーションは、冷蔵庫の食材管理と GPT による献立提案という 2 つの主要機能を提供します。

技術スタック：
- React 18 + TypeScript
- Vite (ビルドツール)
- React Router (ルーティング)
- Zustand (状態管理)
- Tailwind CSS (スタイリング)
- Vitest + React Testing Library (テスト)
- Biome (Lint/Format)

## Architecture

### アーキテクチャパターン

フロントエンドは、レイヤードアーキテクチャを採用します：

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│  (Pages, Components, Layouts)           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Application Layer               │
│  (Hooks, State Management)              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Service Layer                   │
│  (API Client, Data Transformation)      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         External Services               │
│  (Backend API, LocalStorage)            │
└─────────────────────────────────────────┘
```

### ディレクトリ構成

```
src/
├── main.tsx                 # アプリケーションエントリポイント
├── App.tsx                  # ルートコンポーネント
├── pages/                   # ページコンポーネント
│   ├── Home.tsx            # 食材管理ページ
│   ├── Result.tsx          # 献立提案結果ページ
│   ├── Settings.tsx        # 設定ページ
│   └── Meals.tsx           # お気に入り献立ページ
├── components/              # 再利用可能なコンポーネント
│   ├── Header.tsx
│   ├── Footer.tsx
│   ├── MealCard.tsx        # 献立カード
│   ├── IngredientList.tsx  # 食材リスト
│   ├── IngredientForm.tsx  # 食材入力フォーム
│   └── LoadingSpinner.tsx  # ローディング表示
├── layouts/
│   └── DefaultLayout.tsx   # 共通レイアウト
├── hooks/                   # カスタムフック
│   ├── useIngredients.ts   # 食材管理ロジック
│   └── useRecipes.ts       # 献立提案ロジック
├── stores/                  # Zustand ストア
│   ├── ingredientStore.ts  # 食材状態管理
│   └── recipeStore.ts      # 献立状態管理
├── services/                # API クライアント
│   ├── api.ts              # API 基本設定
│   ├── ingredientApi.ts    # 食材 API
│   └── recipeApi.ts        # 献立 API
├── types/                   # TypeScript 型定義
│   └── index.ts
├── utils/                   # ユーティリティ関数
│   ├── validation.ts       # バリデーション
│   └── formatters.ts       # データフォーマット
├── styles/
│   └── globals.css         # グローバルスタイル
└── assets/                  # 静的アセット
```

## Components and Interfaces

### 主要コンポーネント

#### 1. Home Page
食材管理のメインページです。

**責務:**
- 食材一覧の表示
- 食材の追加・編集・削除
- 献立提案のリクエスト

**Props:** なし

**State:**
- `ingredients`: 食材リスト (Zustand store から取得)
- `isLoading`: ローディング状態
- `error`: エラーメッセージ

**主要メソッド:**
- `handleAddIngredient()`: 食材追加処理
- `handleEditIngredient(id)`: 食材編集処理
- `handleDeleteIngredient(id)`: 食材削除処理
- `handleRequestSuggestion()`: 献立提案リクエスト

#### 2. IngredientForm Component
食材の入力フォームコンポーネントです。

**Props:**
```typescript
interface IngredientFormProps {
  ingredient?: Ingredient;  // 編集時は既存データ
  onSubmit: (data: IngredientInput) => void;
  onCancel: () => void;
}
```

**バリデーション:**
- `name`: 必須、1-255文字
- `quantity`: 任意、1-100文字
- `expiration_date`: 任意、YYYY-MM-DD形式

#### 3. MealCard Component
献立を表示するカードコンポーネントです。

**Props:**
```typescript
interface MealCardProps {
  recipe: Recipe;
  onSaveFavorite: (recipe: Recipe) => void;
}
```

**表示内容:**
- 料理名
- 調理手順
- 不足している食材

#### 4. Result Page
献立提案結果を表示するページです。

**責務:**
- 提案された献立の一覧表示
- お気に入りへの保存

**Props:** なし

**State:**
- `recipes`: 献立リスト (Zustand store から取得)

### カスタムフック

#### useIngredients
食材管理のロジックをカプセル化します。

```typescript
interface UseIngredientsReturn {
  ingredients: Ingredient[];
  isLoading: boolean;
  error: string | null;
  fetchIngredients: () => Promise<void>;
  addIngredient: (data: IngredientInput) => Promise<void>;
  updateIngredient: (id: number, data: IngredientInput) => Promise<void>;
  deleteIngredient: (id: number) => Promise<void>;
}
```

#### useRecipes
献立提案のロジックをカプセル化します。

```typescript
interface UseRecipesReturn {
  recipes: Recipe[];
  isLoading: boolean;
  error: string | null;
  requestSuggestion: () => Promise<void>;
  saveFavorite: (recipe: Recipe) => void;
  getFavorites: () => Recipe[];
}
```

## Data Models

### Ingredient (食材)

```typescript
interface Ingredient {
  id: number;
  name: string;
  quantity?: string;
  expiration_date?: string;  // YYYY-MM-DD
}

interface IngredientInput {
  name: string;
  quantity?: string;
  expiration_date?: string;
}
```

### Recipe (献立)

```typescript
interface Recipe {
  name: string;
  steps: string[];
  missing_items: string[];
}

interface RecipeSuggestionResponse {
  suggestions: Recipe[];
}
```

### API Response Types

```typescript
interface ApiError {
  message: string;
  code?: string;
}

interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}
```

## Service Layer

### API Client 設計

#### Base API Configuration

```typescript
// services/api.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const apiClient = {
  get: async <T>(endpoint: string): Promise<T> => { /* ... */ },
  post: async <T>(endpoint: string, data: unknown): Promise<T> => { /* ... */ },
  put: async <T>(endpoint: string, data: unknown): Promise<T> => { /* ... */ },
  delete: async <T>(endpoint: string): Promise<T> => { /* ... */ },
};
```

#### Ingredient API

```typescript
// services/ingredientApi.ts
export const ingredientApi = {
  getAll: () => apiClient.get<Ingredient[]>('/ingredients'),
  create: (data: IngredientInput) => apiClient.post<Ingredient>('/ingredients', data),
  update: (id: number, data: IngredientInput) => 
    apiClient.put<Ingredient>(`/ingredients/${id}`, data),
  delete: (id: number) => apiClient.delete<void>(`/ingredients/${id}`),
};
```

#### Recipe API

```typescript
// services/recipeApi.ts
export const recipeApi = {
  requestSuggestion: () => 
    apiClient.post<RecipeSuggestionResponse>('/recipes/suggestion', {}),
};
```

## State Management

### Zustand Store 設計

#### Ingredient Store

```typescript
interface IngredientState {
  ingredients: Ingredient[];
  isLoading: boolean;
  error: string | null;
  setIngredients: (ingredients: Ingredient[]) => void;
  addIngredient: (ingredient: Ingredient) => void;
  updateIngredient: (id: number, ingredient: Ingredient) => void;
  removeIngredient: (id: number) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
}
```

#### Recipe Store

```typescript
interface RecipeState {
  recipes: Recipe[];
  favorites: Recipe[];
  isLoading: boolean;
  error: string | null;
  setRecipes: (recipes: Recipe[]) => void;
  addFavorite: (recipe: Recipe) => void;
  removeFavorite: (recipeName: string) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
}
```

### LocalStorage 永続化

お気に入り献立は localStorage に保存します：

```typescript
const FAVORITES_KEY = 'dinner-decider-favorites';

const saveFavoritesToStorage = (favorites: Recipe[]) => {
  localStorage.setItem(FAVORITES_KEY, JSON.stringify(favorites));
};

const loadFavoritesFromStorage = (): Recipe[] => {
  const stored = localStorage.getItem(FAVORITES_KEY);
  return stored ? JSON.parse(stored) : [];
};
```

## Error Handling

### エラー処理戦略

1. **API エラー**: HTTP ステータスコードに基づいてエラーメッセージを表示
2. **ネットワークエラー**: 接続エラー時は再試行を促すメッセージを表示
3. **バリデーションエラー**: フォーム入力時にリアルタイムでエラーを表示
4. **予期しないエラー**: 汎用エラーメッセージを表示し、コンソールにログ出力

### エラーメッセージ設計

```typescript
const ERROR_MESSAGES = {
  NETWORK_ERROR: 'ネットワークエラーが発生しました。接続を確認してください。',
  SERVER_ERROR: 'サーバーエラーが発生しました。しばらくしてから再試行してください。',
  NOT_FOUND: '指定されたリソースが見つかりませんでした。',
  VALIDATION_ERROR: '入力内容に誤りがあります。',
  UNKNOWN_ERROR: '予期しないエラーが発生しました。',
};
```

### エラーバウンダリ

React Error Boundary を使用して、予期しないエラーをキャッチします：

```typescript
class ErrorBoundary extends React.Component<Props, State> {
  // エラー発生時にフォールバック UI を表示
}
```

## Testing Strategy

### テストレベル

#### 1. Unit Tests (単体テスト)
- **対象**: ユーティリティ関数、カスタムフック、ストア
- **ツール**: Vitest
- **カバレッジ目標**: 80%以上

**例:**
```typescript
// utils/validation.test.ts
describe('validateIngredient', () => {
  it('should return error when name is empty', () => {
    const result = validateIngredient({ name: '' });
    expect(result.errors.name).toBeDefined();
  });
});
```

#### 2. Component Tests (コンポーネントテスト)
- **対象**: 個別コンポーネント
- **ツール**: Vitest + React Testing Library
- **重点**: ユーザーインタラクション、条件付きレンダリング

**例:**
```typescript
// components/IngredientForm.test.tsx
describe('IngredientForm', () => {
  it('should call onSubmit with form data', async () => {
    const onSubmit = vi.fn();
    render(<IngredientForm onSubmit={onSubmit} onCancel={() => {}} />);
    
    await userEvent.type(screen.getByLabelText('食材名'), 'にんじん');
    await userEvent.click(screen.getByRole('button', { name: '追加' }));
    
    expect(onSubmit).toHaveBeenCalledWith({ name: 'にんじん' });
  });
});
```

#### 3. Integration Tests (統合テスト)
- **対象**: ページレベルのフロー
- **ツール**: Vitest + React Testing Library + MSW (Mock Service Worker)
- **重点**: API 連携、状態管理、ページ遷移

**例:**
```typescript
// pages/Home.test.tsx
describe('Home Page Integration', () => {
  it('should display ingredients and allow adding new one', async () => {
    // MSW でモック API を設定
    server.use(
      rest.get('/ingredients', (req, res, ctx) => {
        return res(ctx.json([{ id: 1, name: 'にんじん' }]));
      })
    );
    
    render(<Home />);
    
    // 既存の食材が表示される
    expect(await screen.findByText('にんじん')).toBeInTheDocument();
    
    // 新しい食材を追加
    await userEvent.click(screen.getByRole('button', { name: '追加' }));
    // ...
  });
});
```

### テスト実行戦略

- **開発時**: `vitest --watch` でファイル変更時に自動実行
- **CI**: PR 作成時に全テストを実行
- **カバレッジ**: `vitest --coverage` でカバレッジレポートを生成

## Accessibility

### WCAG AA 準拠

#### 1. キーボードナビゲーション
- すべてのインタラクティブ要素は Tab キーでフォーカス可能
- Enter/Space キーでボタンを操作可能
- Esc キーでモーダル/ダイアログを閉じる

#### 2. ARIA 属性
```typescript
// 例: ローディング状態
<div role="status" aria-live="polite" aria-busy={isLoading}>
  {isLoading ? 'Loading...' : 'Content'}
</div>

// 例: エラーメッセージ
<div role="alert" aria-live="assertive">
  {error}
</div>
```

#### 3. カラーコントラスト
- テキストと背景のコントラスト比: 最低 4.5:1
- 大きなテキスト (18pt以上): 最低 3:1
- Tailwind CSS のデフォルトカラーパレットを使用

#### 4. フォーカスインジケータ
```css
/* globals.css */
*:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}
```

### スクリーンリーダー対応

- 画像には適切な `alt` 属性を設定
- フォーム要素には `<label>` を関連付け
- ページタイトルは `<h1>` で明示

## Performance Optimization

### 1. コード分割
```typescript
// App.tsx
const Home = lazy(() => import('./pages/Home'));
const Result = lazy(() => import('./pages/Result'));
const Settings = lazy(() => import('./pages/Settings'));
```

### 2. バンドルサイズ最適化
- Tree shaking を有効化
- 不要な依存関係を削除
- 画像は WebP 形式を使用

### 3. Service Worker によるキャッシュ
```typescript
// vite.config.ts で PWA プラグインを設定
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    VitePWA({
      registerType: 'autoUpdate',
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg}'],
      },
    }),
  ],
});
```

### 4. レンダリング最適化
- `React.memo` で不要な再レンダリングを防止
- `useMemo` / `useCallback` で計算結果をメモ化

## Responsive Design

### ブレークポイント

Tailwind CSS のデフォルトブレークポイントを使用：

```
sm: 640px   // スマートフォン (横向き)
md: 768px   // タブレット
lg: 1024px  // デスクトップ
xl: 1280px  // 大画面デスクトップ
```

### モバイルファースト設計

```typescript
// 例: レスポンシブグリッド
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* カード */}
</div>
```

### タッチターゲットサイズ

- ボタン: 最低 44x44px
- リンク: 最低 44x44px
- フォーム入力: 最低 44px の高さ

## Build and Deployment

### ビルド設定

```typescript
// vite.config.ts
export default defineConfig({
  plugins: [react()],
  build: {
    target: 'es2015',
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
        },
      },
    },
  },
});
```

### 環境変数

```
VITE_API_BASE_URL=http://localhost:8080
```

### デプロイ先候補

1. **Vercel**: 推奨。自動デプロイ、プレビュー環境、高速 CDN
2. **Netlify**: Vercel の代替
3. **GitHub Pages**: 静的サイトホスティング

### CI/CD パイプライン

```yaml
# .github/workflows/ci.yml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npm run lint
      - run: npm run type-check
      - run: npm run test
      - run: npm run build
```

## Security Considerations

### 1. XSS 対策
- React のデフォルトエスケープを活用
- `dangerouslySetInnerHTML` は使用しない

### 2. CSRF 対策
- バックエンドで CSRF トークンを実装
- フロントエンドでトークンをヘッダーに含める

### 3. 環境変数の管理
- API キーなどの機密情報は環境変数で管理
- `.env` ファイルは `.gitignore` に追加

### 4. HTTPS の使用
- 本番環境では必ず HTTPS を使用
- Mixed Content を避ける
