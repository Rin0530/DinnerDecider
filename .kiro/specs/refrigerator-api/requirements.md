# Requirements Document

## Introduction

冷蔵庫の中身をMySQLに記録し、LLM（Ollama）に今日の夕飯を考えてもらうWebアプリケーションのバックエンドAPIシステム。ユーザーは食材の登録・更新・削除を行い、現在の食材リストを基にAIから献立提案を受け取ることができます。

## Glossary

- **System**: 冷蔵庫管理バックエンドAPIシステム
- **Client**: Webフロントエンドまたはモバイルアプリケーション
- **Ingredient**: 冷蔵庫に保管される食材（名前、数量、賞味期限を含む）
- **Recipe Suggestion**: LLMが生成する献立提案（料理名、手順、不足食材を含む）
- **LLM Service**: Ollama APIを使用した外部言語モデルサービス
- **Database**: MySQL データベース

## Requirements

### Requirement 1

**User Story:** ユーザーとして、冷蔵庫に新しい食材を追加したい。そうすることで、現在の在庫を正確に管理できる。

#### Acceptance Criteria

1. WHEN Client sends a POST request to `/ingredients` with valid ingredient data, THE System SHALL store the ingredient in the Database and return the created ingredient with a unique ID and HTTP status 201
2. IF Client sends a POST request to `/ingredients` without a required name field, THEN THE System SHALL reject the request and return an error message with HTTP status 400
3. THE System SHALL accept optional quantity and expiration_date fields when creating an ingredient
4. THE System SHALL automatically set created_at and updated_at timestamps when storing a new ingredient

### Requirement 2

**User Story:** ユーザーとして、冷蔵庫に入っている全ての食材を確認したい。そうすることで、何が使えるか把握できる。

#### Acceptance Criteria

1. WHEN Client sends a GET request to `/ingredients`, THE System SHALL retrieve all ingredients from the Database and return them as a JSON array with HTTP status 200
2. IF the Database contains no ingredients, THE System SHALL return an empty array with HTTP status 200
3. THE System SHALL include id, name, quantity, expiration_date, created_at, and updated_at fields for each ingredient in the response

### Requirement 3

**User Story:** ユーザーとして、食材の情報（数量など）を更新したい。そうすることで、使用した分を反映できる。

#### Acceptance Criteria

1. WHEN Client sends a PUT request to `/ingredients/:id` with valid update data, THE System SHALL update the specified ingredient in the Database and return the updated ingredient with HTTP status 200
2. IF Client sends a PUT request to `/ingredients/:id` with a non-existent ID, THEN THE System SHALL return an error message with HTTP status 404
3. THE System SHALL update the updated_at timestamp when modifying an ingredient
4. THE System SHALL allow partial updates of ingredient fields

### Requirement 4

**User Story:** ユーザーとして、使い切った食材を削除したい。そうすることで、リストを最新の状態に保てる。

#### Acceptance Criteria

1. WHEN Client sends a DELETE request to `/ingredients/:id` with a valid ingredient ID, THE System SHALL remove the ingredient from the Database and return HTTP status 204
2. IF Client sends a DELETE request to `/ingredients/:id` with a non-existent ID, THEN THE System SHALL return an error message with HTTP status 404

### Requirement 5

**User Story:** ユーザーとして、現在の食材を使った献立をAIに提案してもらいたい。そうすることで、夕食のメニューを簡単に決められる。

#### Acceptance Criteria

1. WHEN Client sends a POST request to `/recipes/suggestion`, THE System SHALL retrieve all current ingredients from the Database, format them into a prompt, send the prompt to the LLM Service, and return the recipe suggestions with HTTP status 200
2. THE System SHALL format the ingredient list as a comma-separated text string including name and quantity for each ingredient
3. THE System SHALL request the LLM Service to return recipe suggestions in JSON format containing recipe name, cooking steps, and missing ingredients
4. IF the LLM Service returns an error or times out, THEN THE System SHALL return an error message to the Client with HTTP status 503
5. IF the Database contains no ingredients, THE System SHALL inform the LLM Service that no ingredients are available and request general dinner suggestions

### Requirement 6

**User Story:** システム管理者として、データベース接続やAPI設定を環境ごとに変更したい。そうすることで、開発環境と本番環境を適切に分離できる。

#### Acceptance Criteria

1. THE System SHALL read configuration from environment variables or a configuration file at startup
2. THE System SHALL support configuration of Database connection parameters including host, port, username, password, and database name
3. THE System SHALL support configuration of LLM Service endpoint URL and API credentials
4. THE System SHALL support configuration of server port and host address

### Requirement 7

**User Story:** 開発者として、APIエラーを適切にハンドリングしたい。そうすることで、クライアントが問題を理解し対処できる。

#### Acceptance Criteria

1. WHEN an error occurs during request processing, THE System SHALL return a JSON response containing an error message and appropriate HTTP status code
2. THE System SHALL return HTTP status 400 for client input validation errors
3. THE System SHALL return HTTP status 404 for resource not found errors
4. THE System SHALL return HTTP status 500 for internal server errors
5. THE System SHALL return HTTP status 503 for external service unavailability errors
6. THE System SHALL log all errors with sufficient detail for debugging purposes
