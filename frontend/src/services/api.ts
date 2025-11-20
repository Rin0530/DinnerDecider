import type { ApiError } from '../types';

// APIのベースURL（環境変数から取得、デフォルトはlocalhost）
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api';

/**
 * APIクライアントクラス
 * バックエンドAPIとの通信を管理
 */
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  /**
   * HTTPリクエストを実行する共通メソッド
   */
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);

      // レスポンスがエラーの場合
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          message: `HTTP error! status: ${response.status}`,
        }));

        const error: ApiError = {
          error: errorData.error,
          message: errorData.message || `HTTP error! status: ${response.status}`,
          code: errorData.code || `HTTP_${response.status}`,
        };

        throw error;
      }

      // 204 No Contentの場合
      if (response.status === 204) {
        return undefined as T;
      }

      return await response.json();
    } catch (error) {
      if ((error as ApiError).message) {
        throw error;
      }

      // ネットワークエラーまたはその他のfetchエラー
      const apiError: ApiError = {
        message: 'ネットワークエラーが発生しました。接続を確認してください。',
        code: 'NETWORK_ERROR',
      };
      throw apiError;
    }
  }

  /**
   * GETリクエストを実行
   */
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'GET',
    });
  }

  /**
   * POSTリクエストを実行
   */
  async post<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  /**
   * PUTリクエストを実行
   */
  async put<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  /**
   * DELETEリクエストを実行
   */
  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
    });
  }
}

// APIクライアントのインスタンスをエクスポート
export const apiClient = new ApiClient(API_BASE_URL);
