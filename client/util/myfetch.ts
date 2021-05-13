const defaultHeaders = {
  'Content-Type': 'application/json',
} as const;

const UNHANDLED_ERROR_MESSAGE = 'サーバーエラーです。しばらく経ってからやり直してください。' as const;

/**
 * 簡素なオレオレfetchAPIラッパー
 * Content-Typeはapplication/json前提
 * AuthorizationにlocalStrageのトークンを入れる
 * エラー時のレスポンスボディはプレーンテキスト前提
 *
 * @returns response.json()
 * @throws fetchがfetchに失敗したError
 * @throws response.json()がパースに失敗したError
 */
export const myfetch = (method: 'POST' | 'GET' | 'PUT' | 'DELETE' | 'PATCH', url: string, body?: any): Promise<any> => {
  const jwtToken = localStorage.getItem('MyToken');
  const headers = jwtToken
    ? {
        ...defaultHeaders,
        Authorization: `Bearer ${jwtToken}`,
      }
    : defaultHeaders;

  return new Promise((resolve, reject) => {
    fetch(url, {
      headers,
      method,
      body: body ? JSON.stringify(body) : undefined,
    }).then((response) => {
      if (!response.ok) {
        response.text().then((message) => {
          reject(new Error(message));
        });

        return;
      }

      if (response.status === 204) {
        resolve(undefined);
      }

      response
        .json()
        .then((json) => {
          resolve(json);
        })
        .catch(() => {
          // jsonパースエラー。想定していないレスポンスが帰ってきている。
          reject(new Error(UNHANDLED_ERROR_MESSAGE));
        });
    });
  });
};
