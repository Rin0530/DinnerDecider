# Implementation Plan

- [x] 1. プロジェクト構造とコア設定のセットアップ
  - backendディレクトリ内にGoモジュールを初期化
  - 必要な依存関係をインストール（Gin, sqlx, Viper, validator, MySQL driver）
  - ディレクトリ構造を作成（cmd/api, internal/handler, internal/usecase, internal/repository, internal/service, internal/domain, pkg/config, pkg/logger）
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [x] 2. 設定管理とロギングの実装
  - [x] 2.1 設定ファイル構造とViperによる設定読み込みを実装
    - config.yamlのサンプルファイルを作成
    - pkg/config パッケージで設定構造体と読み込みロジックを実装
    - 環境変数オーバーライドをサポート
    - _Requirements: 6.1, 6.2, 6.3, 6.4_
  - [x] 2.2 ロギング機能を実装
    - pkg/logger パッケージでログ初期化とヘルパー関数を実装
    - 構造化ログ（JSON形式）をサポート
    - ログレベル設定をサポート
    - _Requirements: 7.6_

- [x] 3. ドメインモデルの定義
  - internal/domain パッケージにIngredient、RecipeSuggestion、RecipeResponse構造体を定義
  - 日付フィールドのJSON marshaling/unmarshalingを適切に処理
  - _Requirements: 1.1, 1.3, 2.3, 5.3_

- [x] 4. データベース接続とRepository層の実装
  - [x] 4.1 データベース接続ユーティリティを実装
    - sqlxとMySQL driverの依存関係を追加
    - pkg/database パッケージを作成
    - MySQL接続プールを初期化する関数を実装
    - 接続エラーハンドリングを実装
    - _Requirements: 6.2_
  - [x] 4.2 IngredientRepositoryインターフェースと実装を作成
    - internal/repository/ingredient.go にインターフェースを定義
    - internal/repository/ingredient_impl.go に実装を作成
    - Create、GetAll、GetByID、Update、Delete メソッドを実装
    - SQLクエリを実装（sqlxを使用）
    - エラーハンドリングを実装
    - _Requirements: 1.1, 1.4, 2.1, 2.3, 3.1, 3.3, 4.1_
  - [x] 4.3 Repository層のユニットテストを作成






    - sqlmockを使用してCRUD操作をテスト
    - 正常系と異常系のテストケースを実装
    - _Requirements: 1.1, 2.1, 3.1, 4.1_

- [x] 5. 外部サービス層（Ollama統合）の実装





  - [x] 5.1 OllamaServiceインターフェースと実装を作成


    - internal/service/ollama.go にインターフェースを定義
    - internal/service/ollama_impl.go に実装を作成
    - HTTPクライアントを使用してOllama APIと通信
    - プロンプトテンプレートを実装
    - 食材リストのフォーマット処理を実装
    - _Requirements: 5.1, 5.2_
  - [x] 5.2 レスポンスパースとエラーハンドリングを実装

    - JSON形式のレスポンスをRecipeResponse構造体にパース
    - タイムアウト処理を実装（30秒）
    - API通信エラーのハンドリングを実装
    - _Requirements: 5.3, 5.4_
  - [x] 5.3 OllamaServiceのユニットテストを作成











    - モックHTTPサーバーを使用してAPI通信をテスト
    - タイムアウトとエラーケースをテスト
    - _Requirements: 5.1, 5.4_

- [x] 6. Usecase層の実装





  - [x] 6.1 リクエスト/レスポンスDTOを定義


    - internal/usecase/dto.go を作成
    - CreateIngredientRequest、UpdateIngredientRequest、ErrorResponse構造体を定義
    - バリデーションタグを追加
    - _Requirements: 1.1, 1.2, 2.1, 3.1_
  - [x] 6.2 IngredientUsecaseインターフェースと実装を作成


    - internal/usecase/ingredient.go にインターフェースを定義
    - internal/usecase/ingredient_impl.go に実装を作成
    - CreateIngredient、GetAllIngredients、UpdateIngredient、DeleteIngredient メソッドを実装
    - リクエストDTOからドメインモデルへの変換を実装
    - バリデーションロジックを実装
    - _Requirements: 1.1, 1.2, 2.1, 3.1, 3.4, 4.1_
  - [x] 6.3 RecipeUsecaseインターフェースと実装を作成


    - internal/usecase/recipe.go にインターフェースを定義
    - internal/usecase/recipe_impl.go に実装を作成
    - GetRecipeSuggestion メソッドを実装
    - Repository経由で食材を取得
    - OllamaService経由でレシピ提案を取得
    - 食材が空の場合の処理を実装
    --_Requirements: 5.1, 5.5_

  - [x] 6.4 Usecase層のユニットテストを作成





    - モックRepositoryとモックServiceを使用
    - ビジネスロジックとエラーハンドリングをテスト
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1_

- [x] 7. Handler層とルーティングの実装





  - [x] 7.1 エラーレスポンスヘルパーを実装


    - internal/handler/error.go を作成
    - HTTPステータスコードマッピングを実装
    - 統一されたエラーレスポンス形式を実装
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_
  - [x] 7.2 IngredientHandlerを実装


    - internal/handler/ingredient.go を作成
    - POST /ingredients エンドポイントを実装
    - GET /ingredients エンドポイントを実装
    - PUT /ingredients/:id エンドポイントを実装
    - DELETE /ingredients/:id エンドポイントを実装
    - リクエストバリデーションを実装
    - _Requirements: 1.1, 1.2, 2.1, 2.2, 3.1, 3.2, 4.1, 4.2_
  - [x] 7.3 RecipeHandlerを実装


    - internal/handler/recipe.go を作成
    - POST /recipes/suggestion エンドポイントを実装
    - エラーハンドリングを実装
    - _Requirements: 5.1, 5.4_
  - [x] 7.4 ヘルスチェックエンドポイントを実装


    - internal/handler/health.go を作成
    - GET /health エンドポイントを実装
    - GET /health/db エンドポイントを実装（データベース接続確認）



    - GET /health/ollama エンドポイントを実装（Ollama API接続確認）
  - [ ] 7.5 Handler層のユニットテストを作成


    - httptestを使用してHTTPリクエスト/レスポンスをテスト
    - バリデーションエラーケースをテスト
    - _Requirements: 1.1, 1.2, 2.1, 3.1, 4.1, 5.1_

- [x] 8. アプリケーションエントリポイントの実装





  - cmd/api/main.go を実装
  - 設定読み込み、ロガー初期化、データベース接続を実行
  - 依存性注入（Repository、Service、Usecase、Handlerの初期化）を実装
  - Ginルーターのセットアップとサーバー起動を実装
  - グレースフルシャットダウンを実装
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [x] 9. データベースマイグレーション





  - ingredientsテーブルのCREATE文を含むSQLファイルを作成
  - マイグレーション実行手順をREADMEに記載
  - _Requirements: 1.1, 1.4, 2.3, 3.3_

- [x] 10. 統合テストとドキュメント





  - [x] 10.1 統合テストを作成



    - testcontainers-goを使用してテスト用MySQLコンテナを起動
    - エンドツーエンドのAPIフローをテスト
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1_

  - [x] 10.2 README.mdを更新

    - プロジェクト概要、セットアップ手順、API仕様を記載
    - 環境変数とconfig.yamlの説明を記載
    - 開発・実行手順を記載
