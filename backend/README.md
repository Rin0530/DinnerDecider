# Refrigerator API - Backend

冷蔵庫管理バックエンドAPIシステム

冷蔵庫の中身をMySQLに記録し、LLM（Ollama）を使って今日の夕飯の献立を提案するWebアプリケーションのバックエンドAPIです。

## 目次

- [プロジェクト概要](#プロジェクト概要)
- [プロジェクト構造](#プロジェクト構造)
- [技術スタック](#技術スタック)
- [セットアップ](#セットアップ)
- [設定](#設定)
- [実行](#実行)
- [API仕様](#api仕様)
- [テスト](#テスト)
- [トラブルシューティング](#トラブルシューティング)

## プロジェクト概要

このAPIシステムは以下の機能を提供します：

- **食材管理**: 冷蔵庫内の食材の登録、更新、削除、一覧取得
- **レシピ提案**: 現在の食材を基にLLM（Ollama）が献立を提案
- **ヘルスチェック**: アプリケーション、データベース、外部サービスの稼働状態確認

## プロジェクト構造

```
backend/
├── cmd/
│   └── api/              # アプリケーションエントリポイント
│       └── main.go       # メイン関数、サーバー起動
├── internal/
│   ├── handler/          # HTTPハンドラ層（リクエスト/レスポンス処理）
│   ├── usecase/          # ビジネスロジック層
│   ├── repository/       # データアクセス層（データベース操作）
│   ├── service/          # 外部サービス連携層（Ollama API）
│   └── domain/           # ドメインモデル（エンティティ定義）
├── pkg/
│   ├── config/           # 設定管理（Viper）
│   ├── database/         # データベース接続（sqlx）
│   └── logger/           # ロギング（logrus）
├── migrations/           # データベースマイグレーション
│   └── 001_create_ingredients_table.sql
├── integration_test.go   # 統合テスト
├── config.yaml           # 設定ファイル
├── go.mod                # Go モジュール定義
└── README.md             # このファイル
```

## 技術スタック

- **言語**: Go 1.21以上
- **Webフレームワーク**: [Gin](https://github.com/gin-gonic/gin) -
  高速なHTTPルーター
- **データベース**: MySQL 8.0以上 with [sqlx](https://github.com/jmoiron/sqlx)
- **設定管理**: [Viper](https://github.com/spf13/viper) -
  環境変数とYAMLファイルのサポート
- **バリデーション**:
  [go-playground/validator](https://github.com/go-playground/validator)
- **MySQLドライバ**:
  [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- **ロギング**: [logrus](https://github.com/sirupsen/logrus) - 構造化ログ
- **テスト**:
  [testcontainers-go](https://github.com/testcontainers/testcontainers-go) -
  統合テスト用

## アーキテクチャ

このプロジェクトはクリーンアーキテクチャの原則に基づいて設計されています：

- **Handler層**: HTTPリクエスト/レスポンスの処理、バリデーション
- **Usecase層**: ビジネスロジックの実装、トランザクション管理
- **Repository層**: データベースアクセスの抽象化、CRUD操作
- **Service層**: 外部API（Ollama）との通信
- **Domain層**: ドメインモデルの定義

各層は明確に分離されており、依存関係は外側から内側への一方向のみです。

## セットアップ

### 前提条件

以下のソフトウェアがインストールされている必要があります：

- **Go 1.21以上**: [公式サイト](https://golang.org/dl/)からダウンロード
- **MySQL 8.0以上**:
  [公式サイト](https://dev.mysql.com/downloads/mysql/)からダウンロード
- **Ollama**（オプション、レシピ提案機能を使用する場合）:
  [公式サイト](https://ollama.ai/)からダウンロード
- **Docker**（オプション、統合テストを実行する場合）:
  [公式サイト](https://www.docker.com/)からダウンロード

### 依存関係のインストール

プロジェクトディレクトリに移動し、依存関係をダウンロードします：

```bash
cd backend
go mod download
```

### データベースのセットアップ

#### 1. MySQLデータベースの作成

```bash
mysql -u root -p
```

```sql
CREATE DATABASE refrigerator CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'refrigerator_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON refrigerator.* TO 'refrigerator_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 2. マイグレーションの実行

マイグレーションファイルを使用してテーブルを作成します：

```bash
mysql -u refrigerator_user -p refrigerator < migrations/001_create_ingredients_table.sql
```

または、MySQLクライアントから直接実行：

```bash
mysql -u refrigerator_user -p refrigerator
```

```sql
SOURCE migrations/001_create_ingredients_table.sql;
```

#### 3. マイグレーションの確認

テーブルが正しく作成されたことを確認：

```sql
USE refrigerator;
SHOW TABLES;
DESCRIBE ingredients;
```

## 設定

### config.yaml

`config.yaml` ファイルで各種設定を行います。以下は設定項目の説明です：

```yaml
server:
    port: "8080" # APIサーバーのポート番号
    host: "0.0.0.0" # バインドするホストアドレス

database:
    host: "localhost" # MySQLサーバーのホスト
    port: 3306 # MySQLサーバーのポート
    user: "refrigerator_user" # データベースユーザー名
    password: "your_password" # データベースパスワード
    dbname: "refrigerator" # データベース名
    max_open_conns: 25 # 最大オープン接続数
    max_idle_conns: 5 # 最大アイドル接続数

ollama:
    endpoint: "http://localhost:11434" # Ollama APIのエンドポイント
    model: "llama2" # 使用するLLMモデル
    timeout: "30s" # APIタイムアウト時間

logging:
    level: "info" # ログレベル (debug, info, warn, error)
    format: "json" # ログフォーマット (json, text)
```

### 環境変数

環境変数で設定を上書きすることができます。環境変数名は設定項目をアンダースコアで区切り、大文字にします：

```bash
# サーバー設定
export SERVER_PORT=8080
export SERVER_HOST=0.0.0.0

# データベース設定
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USER=refrigerator_user
export DATABASE_PASSWORD=your_password
export DATABASE_DBNAME=refrigerator
export DATABASE_MAX_OPEN_CONNS=25
export DATABASE_MAX_IDLE_CONNS=5

# Ollama設定
export OLLAMA_ENDPOINT=http://localhost:11434
export OLLAMA_MODEL=llama2
export OLLAMA_TIMEOUT=30s

# ロギング設定
export LOGGING_LEVEL=info
export LOGGING_FORMAT=json
```

環境変数は `config.yaml` の設定よりも優先されます。

## 実行

### 開発モード

開発時は `go run` コマンドで直接実行できます：

```bash
go run cmd/api/main.go
```

サーバーが起動すると、以下のようなログが表示されます：

```
INFO[0000] Starting refrigerator API server...
INFO[0000] Server is listening on 0.0.0.0:8080
```

### ビルドして実行

本番環境では、バイナリをビルドして実行します：

```bash
# ビルド
go build -o bin/api cmd/api/main.go

# 実行
./bin/api
```

### グレースフルシャットダウン

サーバーは `SIGINT` (Ctrl+C) または `SIGTERM`
シグナルを受信すると、グレースフルシャットダウンを実行します：

1. 新しいリクエストの受付を停止
2. 処理中のリクエストが完了するまで最大10秒待機
3. データベース接続をクローズ
4. サーバーを終了

## API仕様

### 食材管理エンドポイント

#### POST /api/ingredients

新しい食材を追加します。

**リクエストボディ:**

```json
{
    "name": "にんじん",
    "quantity": "2本",
    "expiration_date": "2025-11-01"
}
```

- `name` (必須): 食材名
- `quantity` (オプション): 数量
- `expiration_date` (オプション): 賞味期限（YYYY-MM-DD形式）

**レスポンス (201 Created):**

```json
{
    "id": 1,
    "name": "にんじん",
    "quantity": "2本",
    "expiration_date": "2025-11-01T00:00:00Z",
    "created_at": "2025-10-25T10:00:00Z",
    "updated_at": "2025-10-25T10:00:00Z"
}
```

**エラーレスポンス (400 Bad Request):**

```json
{
    "error": "validation_error",
    "message": "Name is required"
}
```

#### GET /api/ingredients

全ての食材を取得します。

**レスポンス (200 OK):**

```json
[
    {
        "id": 1,
        "name": "にんじん",
        "quantity": "2本",
        "expiration_date": "2025-11-01T00:00:00Z",
        "created_at": "2025-10-25T10:00:00Z",
        "updated_at": "2025-10-25T10:00:00Z"
    },
    {
        "id": 2,
        "name": "豚バラ肉",
        "quantity": "200g",
        "expiration_date": null,
        "created_at": "2025-10-25T10:05:00Z",
        "updated_at": "2025-10-25T10:05:00Z"
    }
]
```

食材が存在しない場合は空配列 `[]` を返します。

#### PUT /api/ingredients/:id

指定したIDの食材情報を更新します。

**リクエストボディ:**

```json
{
    "name": "にんじん",
    "quantity": "3本",
    "expiration_date": "2025-11-05"
}
```

全てのフィールドはオプションです。指定したフィールドのみが更新されます。

**レスポンス (200 OK):**

```json
{
    "id": 1,
    "name": "にんじん",
    "quantity": "3本",
    "expiration_date": "2025-11-05T00:00:00Z",
    "created_at": "2025-10-25T10:00:00Z",
    "updated_at": "2025-10-25T11:00:00Z"
}
```

**エラーレスポンス (404 Not Found):**

```json
{
    "error": "not_found",
    "message": "Ingredient not found"
}
```

#### DELETE /api/ingredients/:id

指定したIDの食材を削除します。

**レスポンス (204 No Content):**

レスポンスボディなし

**エラーレスポンス (404 Not Found):**

```json
{
    "error": "not_found",
    "message": "Ingredient not found"
}
```

### レシピ提案エンドポイント

#### POST /api/recipes/suggestion

現在の食材を使った献立をLLMに提案してもらいます。

**リクエストボディ:**

なし（現在データベースに登録されている全ての食材を使用）

**レスポンス (200 OK):**

```json
{
    "suggestions": [
        {
            "name": "肉じゃが",
            "steps": [
                "にんじんとじゃがいもを一口大に切る",
                "豚バラ肉を炒める",
                "野菜を加えて炒め、だし汁と調味料を加える",
                "弱火で20分煮込む"
            ],
            "missing_items": ["じゃがいも", "玉ねぎ"]
        },
        {
            "name": "野菜炒め",
            "steps": [
                "にんじんを千切りにする",
                "豚バラ肉を炒める",
                "にんじんを加えて炒める",
                "塩コショウで味付け"
            ],
            "missing_items": []
        }
    ]
}
```

**エラーレスポンス (503 Service Unavailable):**

```json
{
    "error": "service_unavailable",
    "message": "Recipe suggestion service is currently unavailable"
}
```

Ollama APIが利用できない場合やタイムアウトした場合に返されます。

### ヘルスチェックエンドポイント

#### GET /health

アプリケーションの稼働状態を確認します。

**レスポンス (200 OK):**

```json
{
    "status": "ok"
}
```

#### GET /health/db

データベース接続状態を確認します。

**レスポンス (200 OK):**

```json
{
    "status": "ok",
    "database": "connected"
}
```

**エラーレスポンス (503 Service Unavailable):**

```json
{
    "status": "error",
    "database": "disconnected",
    "error": "connection error details"
}
```

#### GET /health/ollama

Ollama API接続状態を確認します。

**レスポンス (200 OK):**

```json
{
    "status": "ok",
    "ollama": "connected"
}
```

**エラーレスポンス (503 Service Unavailable):**

```json
{
    "status": "error",
    "ollama": "disconnected",
    "error": "connection error details"
}
```

## テスト

### ユニットテスト

各層のユニットテストを実行：

```bash
go test ./internal/... -v
```

### 統合テスト

統合テストは testcontainers-go
を使用して、実際のMySQLコンテナを起動してテストします。

**前提条件**: Dockerが起動している必要があります。

```bash
# 統合テストを実行
go test -v -run TestIngredientCRUDFlow -timeout 120s

# 全てのテストを実行
go test -v -timeout 120s

# カバレッジレポートを生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### テストの種類

- **TestIngredientCRUDFlow**: 食材のCRUD操作の完全なフローをテスト
- **TestValidationErrors**: 入力バリデーションエラーのテスト
- **TestHealthEndpoints**: ヘルスチェックエンドポイントのテスト
- **TestEmptyIngredientsFlow**: 食材が存在しない場合の動作テスト

## トラブルシューティング

### データベース接続エラー

**症状**: `Failed to connect to database` エラーが発生する

**解決方法**:

1. MySQLサーバーが起動していることを確認：
   ```bash
   # Linux/Mac
   sudo systemctl status mysql

   # または
   ps aux | grep mysql
   ```

2. `config.yaml` の接続情報が正しいことを確認：
   - ホスト名、ポート番号
   - ユーザー名、パスワード
   - データベース名

3. データベースユーザーの権限を確認：
   ```bash
   mysql -u root -p
   ```
   ```sql
   SHOW GRANTS FOR 'refrigerator_user'@'localhost';
   ```

4. ファイアウォールでポート3306が開いていることを確認

### マイグレーションエラー

**症状**: テーブル作成時にエラーが発生する

**解決方法**:

1. データベースが存在することを確認：
   ```sql
   SHOW DATABASES;
   ```

2. ユーザーに適切な権限があることを確認：
   ```sql
   GRANT ALL PRIVILEGES ON refrigerator.* TO 'refrigerator_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

3. 既にテーブルが存在する場合は削除してから再実行：
   ```sql
   DROP TABLE IF EXISTS ingredients;
   SOURCE migrations/001_create_ingredients_table.sql;
   ```

### Ollama API接続エラー

**症状**: レシピ提案時に503エラーが返される

**解決方法**:

1. Ollamaサーバーが起動していることを確認：
   ```bash
   # Ollamaサーバーを起動
   ollama serve
   ```

2. モデルがダウンロードされていることを確認：
   ```bash
   # モデル一覧を表示
   ollama list

   # モデルをダウンロード（必要な場合）
   ollama pull llama2
   ```

3. `config.yaml` のエンドポイントが正しいことを確認：
   ```yaml
   ollama:
       endpoint: "http://localhost:11434"
   ```

4. ネットワーク接続を確認：
   ```bash
   curl http://localhost:11434/api/tags
   ```

### ポート競合エラー

**症状**: `bind: address already in use` エラーが発生する

**解決方法**:

1. 使用中のポートを確認：
   ```bash
   # Linux/Mac
   lsof -i :8080

   # Windows
   netstat -ano | findstr :8080
   ```

2. 別のポートを使用するように設定を変更：
   ```yaml
   server:
       port: "8081"
   ```

   または環境変数で上書き：
   ```bash
   export SERVER_PORT=8081
   ```

### ログレベルの変更

デバッグ情報を確認したい場合は、ログレベルを `debug` に変更：

```yaml
logging:
    level: "debug"
    format: "json"
```

または環境変数で：

```bash
export LOGGING_LEVEL=debug
```

### 統合テストが失敗する

**症状**: `checked path: $XDG_RUNTIME_DIR` エラーが発生する

**解決方法**:

1. Dockerが起動していることを確認：
   ```bash
   docker ps
   ```

2. Dockerデーモンが実行されていない場合は起動：
   ```bash
   # Linux
   sudo systemctl start docker

   # Mac/Windows
   Docker Desktopを起動
   ```

3. Docker接続を確認：
   ```bash
   docker run hello-world
   ```

## 開発ガイド

### 新しいエンドポイントの追加

1. **Domain層**: `internal/domain/` にエンティティを定義
2. **Repository層**: `internal/repository/`
   にデータアクセスインターフェースと実装を追加
3. **Usecase層**: `internal/usecase/` にビジネスロジックを実装
4. **Handler層**: `internal/handler/` にHTTPハンドラを実装
5. **Router**: `cmd/api/main.go` の `setupRouter` 関数にルートを追加

### コーディング規約

- パッケージ名は小文字の単一単語
- インターフェース名は実装の責務を表す名詞
- エラーは適切にラップして上位層に伝播
- ログには十分なコンテキスト情報を含める
- テストは `_test.go` サフィックスを使用

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
