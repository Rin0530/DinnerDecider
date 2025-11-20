# Design Document

## Overview

冷蔵庫管理バックエンドAPIシステムは、Go言語で構築されるRESTful APIサーバーです。クリーンアーキテクチャの原則に基づき、Handler層、Usecase層、Repository層、そしてExternal Service層に分離された設計を採用します。これにより、各層の責務が明確になり、テスト容易性と保守性が向上します。

主要な技術スタック：
- **Webフレームワーク**: Gin
- **データベース**: MySQL with sqlx
- **LLMサービス**: Ollama API
- **設定管理**: Viper
- **バリデーション**: go-playground/validator

## Architecture

### アーキテクチャ図

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                               │
│                    (Web/Mobile App)                          │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/JSON
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                           │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │  Ingredient  │  │    Recipe    │                        │
│  │   Handler    │  │   Handler    │                        │
│  └──────┬───────┘  └──────┬───────┘                        │
└─────────┼──────────────────┼──────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                      Usecase Layer                           │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │  Ingredient  │  │    Recipe    │                        │
│  │   Usecase    │  │   Usecase    │                        │
│  └──────┬───────┘  └──────┬───────┘                        │
└─────────┼──────────────────┼──────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────────┐  ┌──────────────────────────────────┐
│  Repository Layer   │  │   External Service Layer         │
│  ┌───────────────┐  │  │  ┌────────────────────────────┐  │
│  │  Ingredient   │  │  │  │     Ollama Service         │  │
│  │  Repository   │  │  │  │    (LLM Integration)       │  │
│  └───────┬───────┘  │  │  └────────────┬───────────────┘  │
└──────────┼──────────┘  └───────────────┼──────────────────┘
           │                              │
           ▼                              ▼
    ┌──────────┐                  ┌──────────────┐
    │  MySQL   │                  │ Ollama API   │
    │ Database │                  │   Server     │
    └──────────┘                  └──────────────┘
```

### レイヤー責務

**Handler Layer (internal/handler)**
- HTTPリクエストの受信とレスポンスの返却
- リクエストボディのバリデーション
- HTTPステータスコードの設定
- Usecaseの呼び出し

**Usecase Layer (internal/usecase)**
- ビジネスロジックの実装
- 複数のRepositoryやServiceの調整
- トランザクション管理
- エラーハンドリング

**Repository Layer (internal/repository)**
- データベースアクセスの抽象化
- CRUD操作の実装
- SQLクエリの実行

**External Service Layer (internal/service)**
- 外部API（Ollama）との通信
- プロンプトの構築
- レスポンスのパース

## Components and Interfaces

### Domain Models (internal/domain)

```go
// Ingredient represents a food item in the refrigerator
type Ingredient struct {
    ID             int64      `json:"id"`
    Name           string     `json:"name"`
    Quantity       string     `json:"quantity"`
    ExpirationDate *time.Time `json:"expiration_date"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}

// RecipeSuggestion represents a recipe suggestion from LLM
type RecipeSuggestion struct {
    Name         string   `json:"name"`
    Steps        []string `json:"steps"`
    MissingItems []string `json:"missing_items"`
}

// RecipeResponse represents the response containing multiple suggestions
type RecipeResponse struct {
    Suggestions []RecipeSuggestion `json:"suggestions"`
}
```

### Repository Interface (internal/repository)

```go
type IngredientRepository interface {
    Create(ctx context.Context, ingredient *domain.Ingredient) error
    GetAll(ctx context.Context) ([]*domain.Ingredient, error)
    GetByID(ctx context.Context, id int64) (*domain.Ingredient, error)
    Update(ctx context.Context, ingredient *domain.Ingredient) error
    Delete(ctx context.Context, id int64) error
}
```

### Usecase Interface (internal/usecase)

```go
type IngredientUsecase interface {
    CreateIngredient(ctx context.Context, req CreateIngredientRequest) (*domain.Ingredient, error)
    GetAllIngredients(ctx context.Context) ([]*domain.Ingredient, error)
    UpdateIngredient(ctx context.Context, id int64, req UpdateIngredientRequest) (*domain.Ingredient, error)
    DeleteIngredient(ctx context.Context, id int64) error
}

type RecipeUsecase interface {
    GetRecipeSuggestion(ctx context.Context) (*domain.RecipeResponse, error)
}
```

### External Service Interface (internal/service)

```go
type OllamaService interface {
    GenerateRecipeSuggestion(ctx context.Context, ingredients []*domain.Ingredient) (*domain.RecipeResponse, error)
}
```

### Handler Structures

**Ingredient Handler**
- `POST /ingredients` → CreateIngredient
- `GET /ingredients` → GetAllIngredients
- `PUT /ingredients/:id` → UpdateIngredient
- `DELETE /ingredients/:id` → DeleteIngredient

**Recipe Handler**
- `POST /recipes/suggestion` → GetRecipeSuggestion

## Data Models

### Database Schema

**ingredients テーブル**

```sql
CREATE TABLE ingredients (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity VARCHAR(100),
    expiration_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_expiration_date (expiration_date)
);
```

### Request/Response DTOs

**CreateIngredientRequest**
```go
type CreateIngredientRequest struct {
    Name           string  `json:"name" binding:"required"`
    Quantity       string  `json:"quantity"`
    ExpirationDate *string `json:"expiration_date"` // YYYY-MM-DD format
}
```

**UpdateIngredientRequest**
```go
type UpdateIngredientRequest struct {
    Name           *string `json:"name"`
    Quantity       *string `json:"quantity"`
    ExpirationDate *string `json:"expiration_date"` // YYYY-MM-DD format
}
```

**ErrorResponse**
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
```

## Error Handling

### エラー分類とHTTPステータスコード

| エラータイプ | HTTPステータス | 使用ケース |
|------------|--------------|----------|
| ValidationError | 400 Bad Request | 入力バリデーション失敗 |
| NotFoundError | 404 Not Found | リソースが存在しない |
| InternalError | 500 Internal Server Error | データベースエラー、予期しないエラー |
| ServiceUnavailableError | 503 Service Unavailable | Ollama APIが利用不可 |

### エラーハンドリング戦略

1. **Handler層**: HTTPステータスコードとJSON形式のエラーレスポンスを返却
2. **Usecase層**: ビジネスロジックエラーを適切なエラータイプに変換
3. **Repository層**: データベースエラーをラップして上位層に伝播
4. **Service層**: 外部APIエラーをラップして上位層に伝播

### ロギング

- すべてのエラーはログに記録
- リクエストID（トレースID）を使用してリクエストを追跡
- ログレベル: DEBUG, INFO, WARN, ERROR

## Testing Strategy

### ユニットテスト

**Repository層**
- モックデータベース（sqlmock）を使用
- CRUD操作の正常系・異常系をテスト

**Usecase層**
- モックRepositoryとモックServiceを使用
- ビジネスロジックの検証
- エラーハンドリングの検証

**Handler層**
- httptest パッケージを使用
- リクエスト/レスポンスの検証
- バリデーションエラーの検証

**Service層**
- モックHTTPクライアントを使用
- Ollama APIとの通信ロジックをテスト

### 統合テスト

- テスト用MySQLコンテナ（testcontainers-go）を使用
- 実際のデータベースを使用したエンドツーエンドテスト
- API全体のフローを検証

### テストカバレッジ目標

- 全体: 70%以上
- Usecase層: 80%以上（ビジネスロジックの重要性が高いため）

## Configuration Management

### 設定項目

```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  dbname: "refrigerator"
  max_open_conns: 25
  max_idle_conns: 5

ollama:
  endpoint: "http://localhost:11434"
  model: "llama2"
  timeout: 30s

logging:
  level: "info"
  format: "json"
```

### 環境変数オーバーライド

Viperを使用して、環境変数で設定をオーバーライド可能：
- `SERVER_PORT`
- `DATABASE_HOST`
- `DATABASE_PASSWORD`
- `OLLAMA_ENDPOINT`

## LLM Integration Details

### プロンプト設計

```
あなたはプロの料理人です。以下の食材を使って作れる、美味しくて簡単な夕食の献立を3つ提案してください。
それぞれの献立には、料理名、簡単な作り方、そして不足している食材（もしあれば）を記載してください。
回答は必ずJSON形式で、以下のフォーマットに従ってください。

{
  "suggestions": [
    {
      "name": "料理名",
      "steps": ["手順1", "手順2", "手順3"],
      "missing_items": ["不足している食材1"]
    }
  ]
}

# 利用可能な食材
{{INGREDIENTS_LIST}}
```

### Ollama API呼び出しフロー

1. Repository経由で全食材を取得
2. 食材リストをフォーマット（例: "にんじん(2本), 豚バラ肉(200g)"）
3. プロンプトテンプレートに食材リストを埋め込み
4. Ollama APIに POST リクエスト送信
5. レスポンスをパースしてRecipeResponse構造体に変換
6. エラーハンドリング（タイムアウト、パースエラー等）

### タイムアウトとリトライ

- デフォルトタイムアウト: 30秒
- リトライなし（初回失敗時は503エラーを返却）
- 将来的な拡張: エクスポネンシャルバックオフによるリトライ実装

## Deployment Considerations

### ヘルスチェックエンドポイント

- `GET /health` - アプリケーションの稼働状態
- `GET /health/db` - データベース接続状態
- `GET /health/ollama` - Ollama API接続状態

## Future Enhancements

1. **ユーザー認証**: JWT認証の追加
2. **マルチユーザー対応**: ユーザーごとの冷蔵庫管理
3. **買い物リスト機能**: 不足食材の買い物リスト生成
4. **食材カテゴリ**: 野菜、肉、魚などのカテゴリ分類
5. **賞味期限アラート**: 期限切れ間近の食材通知
6. **レシピ履歴**: 過去の献立提案の保存と再利用
