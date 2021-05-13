package controller

import (
	"too-simple-twitter/server/util/serverError"
	"encoding/json"
	"net/http"
)

// HTTPリクエストにまつわるヘルパー関数たち

// クライアントのホットリロードサーバが走るホスト
const CORS_ORIGIN = "http://localhost:3000"

/*
  開発用
	localhost:3000からのCORSを許可する
*/
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", CORS_ORIGIN)
		next.ServeHTTP(w, r)
	})
}

/*
  開発用
	CORSプリフライトリクエストをハンドル
*/
func CORSPreflight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, UPDATE, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", CORS_ORIGIN)
}

// リクエストボディをパースする
// 参考：https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func parseRequestBody(result interface{}, w http.ResponseWriter, r *http.Request) serverError.ServerError {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(result)
	if err != nil {
		return serverError.NewFatalServerError(err)
	}

	return nil
}

// serverError.ServerErrorをHTTPレスポンスとして返す
func responseError(err serverError.ServerError, w http.ResponseWriter) {
	http.Error(w, err.Error(), err.Code())
}
