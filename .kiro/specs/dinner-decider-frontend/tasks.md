# Implementation Plan

-
  1. [x] プロジェクトのセットアップと基本構成
  - Vite + React + TypeScript
    プロジェクトを初期化し、必要な依存関係をインストールする
  - Tailwind CSS、React Router、Zustand を設定する
  - Biome を設定し、lint と format のルールを定義する
  - ディレクトリ構造を作成する（pages, components, hooks, stores, services,
    types, utils）
  - _Requirements: 1.1, 1.4, 8.1, 9.1, 10.1_

-
  2. [x] 型定義と API クライアントの実装
  - `types/index.ts` に Ingredient、Recipe、API レスポンスの型定義を作成する
  - `services/api.ts` に基本的な API クライアント（get, post, put, delete
    メソッド）を実装する
  - `services/ingredientApi.ts` に食材 API のエンドポイント関数を実装する
  - `services/recipeApi.ts` に献立 API のエンドポイント関数を実装する
  - 環境変数（VITE_API_BASE_URL）を設定する
  - _Requirements: 1.1, 2.2, 4.2, 5.1_

-
  3. [x] 状態管理ストアの実装
  - `stores/ingredientStore.ts` に Zustand ストアを実装する（食材の状態管理）
  - `stores/recipeStore.ts` に Zustand ストアを実装する（献立の状態管理）
  - localStorage との連携機能を実装する（お気に入り献立の永続化）
  - _Requirements: 5.2, 7.1, 7.3_

-
  4. [x] カスタムフックの実装
  - `hooks/useIngredients.ts` を実装する（食材の取得、追加、更新、削除ロジック）
  - `hooks/useRecipes.ts`
    を実装する（献立提案のリクエスト、お気に入り保存ロジック）
  - エラーハンドリングとローディング状態の管理を含める
  - _Requirements: 1.1, 1.3, 2.2, 2.3, 3.2, 3.3, 4.2, 4.3, 5.1, 5.2, 7.1_

-
  5. [x] 共通コンポーネントの実装
  - `components/Header.tsx` を実装する（ナビゲーション）
  - `components/Footer.tsx` を実装する（フッター情報）
  - `components/LoadingSpinner.tsx` を実装する（ローディング表示）
  - `layouts/DefaultLayout.tsx` を実装する（共通レイアウト）
  - アクセシビリティ対応（ARIA 属性、キーボードナビゲーション）を含める
  - _Requirements: 8.1, 8.2, 8.3, 9.1, 9.2_

-
  6. [x] 食材管理コンポーネントの実装
  - `components/IngredientList.tsx` を実装する（食材一覧表示）
  - `components/IngredientForm.tsx` を実装する（食材の追加・編集フォーム）
  - フォームバリデーション機能を実装する（name: 必須、quantity:
    任意、expiration_date: 日付形式）
  - 編集・削除ボタンと確認ダイアログを実装する
  - _Requirements: 1.2, 2.1, 2.4, 3.1, 4.1_

-
  7. [x] Home ページの実装
  - `pages/Home.tsx` を実装する（食材管理のメインページ）
  - useIngredients フックを使用して食材の取得、追加、編集、削除機能を統合する
  - 「献立を提案」ボタンを実装し、useRecipes フックと連携する
  - エラーメッセージとローディング状態の表示を実装する
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 3.1, 3.2, 3.3, 4.1, 4.2, 4.3,
    5.1, 5.4_

-
  8. [x] 献立表示コンポーネントの実装
  - `components/MealCard.tsx` を実装する（献立カード表示）
  - 料理名、調理手順、不足している食材を表示する
  - 「お気に入りに保存」ボタンを実装する
  - _Requirements: 6.2, 7.1, 7.2_

-
  9. [x] Result ページの実装
  - `pages/Result.tsx` を実装する（献立提案結果ページ）
  - recipeStore から献立リストを取得して表示する
  - MealCard コンポーネントを使用して各献立を表示する
  - 献立がない場合のメッセージを表示する
  - _Requirements: 5.3, 6.1, 6.2, 6.3_

-
  10. [x] ルーティングとナビゲーションの実装
  - [x] 10.1 React Router の基本設定を実装する
    - `App.tsx` に BrowserRouter と Routes を設定する
    - Home、Result、Settings、Meals ページへのルートを定義する
    - _Requirements: 5.3, 10.2_
  - [x] 10.2 コード分割（lazy loading）を実装する

    - React.lazy と Suspense を使用して各ページをコード分割する
    - ローディング中のフォールバック UI を設定する
    - _Requirements: 10.2, 10.3_

-
  11. [x] Settings と Meals ページの実装

  - [x] 11.1 Meals ページを実装する
    - `pages/Meals.tsx` を作成する（お気に入り献立一覧ページ）
    - useRecipeStore から favorites を取得して表示する
    - MealCard コンポーネントを再利用して献立を表示する
    - お気に入りから削除する機能を実装する
    - _Requirements: 7.3_
  - [x] 11.2 Settings ページを実装する
    - `pages/Settings.tsx` を作成する（設定ページの基本構造）
    - 将来の拡張に備えた基本レイアウトを実装する
    - _Requirements: 該当なし（将来の拡張用）_

-
  12. [x] エラーハンドリングの強化

  - [x] 12.1 エラーメッセージ定数を定義する
    - `utils/errorMessages.ts` を作成する
    - NETWORK_ERROR, SERVER_ERROR などの定数を定義する
    - _Requirements: 1.3, 3.4, 5.5_
  - [x] 12.2 Error Boundary コンポーネントを実装する
    - `components/ErrorBoundary.tsx` を作成する
    - 予期しないエラーをキャッチしてフォールバック UI を表示する
    - App.tsx で ErrorBoundary を適用する
    - _Requirements: 1.3, 5.5_

-
  13. [x] レスポンシブデザインの実装
  - Tailwind CSS を使用してモバイルファーストのスタイリングを適用する
  - ブレークポイント（sm, md, lg, xl）に応じたレイアウト調整を実装する
  - タッチターゲットサイズを最低 44x44px に設定する
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

-
  14. [x] アクセシビリティの基本実装
  - すべてのインタラクティブ要素にキーボードナビゲーションを実装する
  - ARIA 属性（role, aria-live, aria-busy など）を追加する
  - フォーカスインジケータのスタイルを設定する
  - カラーコントラストを WCAG AA 基準に準拠させる
  - _Requirements: 8.1, 8.2, 8.3, 8.4_

-
  15. [x] パフォーマンス最適化
  - [x] 15.1 コンポーネントのメモ化を実装する
    - MealCard、IngredientList、IngredientForm に React.memo を適用する
    - 不要な再レンダリングを防止する
    - _Requirements: 10.1, 10.2_
  - [x] 15.2 フックのメモ化を実装する
    - useIngredients と useRecipes で useMemo と useCallback を適用する
    - 計算結果とコールバック関数をメモ化する
    - _Requirements: 10.1, 10.2_

-
  16. [x] Service Worker とオフライン対応
  - vite-plugin-pwa を設定する
  - Service Worker によるキャッシュ戦略を実装する
  - オフライン時の基本的な閲覧機能を提供する
  - _Requirements: 10.4_

-
  17. [x] TypeScript 構文エラーとコード品質の修正
  - TypeScript の型エラーを修正する（React の型定義、JSX の型定義など）
  - `tsconfig.json` の設定を確認・調整する
  - `pnpm run type-check` (tsc --noEmit) が成功することを確認する
  - `pnpm run lint` (biome check) が成功することを確認する
  - `pnpm run format` (biome format --write) を実行してコードをフォーマットする
  - すべてのコード品質チェックが通ることを確認する
  - _Requirements: 該当なし（コード品質）_

-
  18. [ ] テストの実装

  - [ ]* 18.1 テスト環境のセットアップ

    - MSW (Mock Service Worker) をインストールして設定する
    - テストユーティリティファイルを作成する
    - _Requirements: 該当なし（テストインフラ）_
  - [ ]* 18.2 ユーティリティ関数のユニットテストを作成する

    - `utils/validation.test.ts` を作成してバリデーション関数のテストを実装する
    - _Requirements: 2.4_
  - [ ]* 18.3 カスタムフックのテストを作成する

    - `hooks/useIngredients.test.ts` を作成する
    - `hooks/useRecipes.test.ts` を作成する
    - _Requirements: 1.1, 2.2, 3.2, 4.2, 5.1, 7.1_
  - [ ]* 18.4 コンポーネントのテストを作成する

    - `components/IngredientForm.test.tsx` を作成する
    - `components/MealCard.test.tsx` を作成する
    - `components/LoadingSpinner.test.tsx` を作成する
    - _Requirements: 2.1, 6.2_
  - [ ]* 18.5 統合テストを作成する

    - `pages/Home.test.tsx` を作成する（MSW でモック API を使用）
    - `pages/Result.test.tsx` を作成する
    - _Requirements: 1.1, 1.2, 2.2, 5.1, 6.1_

-
  19. [x] ビルド設定と環境変数の最終調整
  - vite.config.ts でビルド最適化を設定する（manualChunks など）
  - .env.example ファイルを作成する
  - README.md にセットアップ手順を記載する
  - _Requirements: 10.1, 10.2_

-
  20. [ ]* CI/CD パイプラインの設定
  - [ ]* 20.1 GitHub Actions ワークフローを作成する
    - `.github/workflows/ci.yml` ファイルを作成する
    - pnpm のセットアップステップを追加する
    - lint (biome check), type-check (tsc --noEmit), test (vitest --run), build (vite build) のステップを定義する
    - PR 作成時とプッシュ時にワークフローを実行する設定を追加する
    - _Requirements: 10.1_
  - [ ]* 20.2 自動デプロイを設定する
    - Vercel または Netlify のデプロイ設定ファイルを作成する
    - 環境変数の設定手順をドキュメント化する
    - main ブランチへのマージ後に自動デプロイされるように設定する
    - _Requirements: 10.1_
