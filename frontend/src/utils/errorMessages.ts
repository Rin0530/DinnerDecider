/**
 * エラーメッセージ定数
 * アプリケーション全体で使用するエラーメッセージを定義
 */

export const ERROR_MESSAGES = {
  // ネットワークエラー
  NETWORK_ERROR: 'ネットワークエラーが発生しました。接続を確認してください。',

  // サーバーエラー
  SERVER_ERROR: 'サーバーエラーが発生しました。しばらくしてから再試行してください。',

  // リソースが見つからない
  NOT_FOUND: '指定されたリソースが見つかりませんでした。',

  // バリデーションエラー
  VALIDATION_ERROR: '入力内容に誤りがあります。',

  // 予期しないエラー
  UNKNOWN_ERROR: '予期しないエラーが発生しました。',

  // 認証エラー
  UNAUTHORIZED: '認証に失敗しました。再度ログインしてください。',

  // 権限エラー
  FORBIDDEN: 'この操作を実行する権限がありません。',

  // タイムアウトエラー
  TIMEOUT_ERROR: 'リクエストがタイムアウトしました。再試行してください。',

  // データ取得エラー
  FETCH_ERROR: 'データの取得に失敗しました。',

  // データ保存エラー
  SAVE_ERROR: 'データの保存に失敗しました。',

  // データ更新エラー
  UPDATE_ERROR: 'データの更新に失敗しました。',

  // データ削除エラー
  DELETE_ERROR: 'データの削除に失敗しました。',
} as const;

/**
 * HTTPステータスコードに基づいてエラーメッセージを取得
 * @param status HTTPステータスコード
 * @returns エラーメッセージ
 */
export const getErrorMessageByStatus = (status: number): string => {
  switch (status) {
    case 400:
      return ERROR_MESSAGES.VALIDATION_ERROR;
    case 401:
      return ERROR_MESSAGES.UNAUTHORIZED;
    case 403:
      return ERROR_MESSAGES.FORBIDDEN;
    case 404:
      return ERROR_MESSAGES.NOT_FOUND;
    case 408:
      return ERROR_MESSAGES.TIMEOUT_ERROR;
    case 500:
    case 502:
    case 503:
    case 504:
      return ERROR_MESSAGES.SERVER_ERROR;
    default:
      return ERROR_MESSAGES.UNKNOWN_ERROR;
  }
};

/**
 * エラーオブジェクトからユーザーフレンドリーなメッセージを取得
 * @param error エラーオブジェクト
 * @returns エラーメッセージ
 */
export const getErrorMessage = (error: unknown): string => {
  if (error instanceof Error) {
    // ネットワークエラーの判定
    if (error.message.includes('fetch') || error.message.includes('network')) {
      return ERROR_MESSAGES.NETWORK_ERROR;
    }
    return error.message;
  }

  if (typeof error === 'string') {
    return error;
  }

  return ERROR_MESSAGES.UNKNOWN_ERROR;
};
