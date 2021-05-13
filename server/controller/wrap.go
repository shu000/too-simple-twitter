package controller

import (
	"too-simple-twitter/server/util/serverError"
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

// 共通処理のコールバック関数
type wrappedFunc func(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError

/*
  全APIの共通処理でコールバック関数をラップする
  共通処理はRDBのトランザクションオープンとエラーハンドル
*/
func wrap(f wrappedFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// DBコネクション
		host := []byte(os.Getenv("POSTGRES_HOST"))
		dbname := []byte(os.Getenv("POSTGRES_DB"))
		user := []byte(os.Getenv("POSTGRES_USER"))
		password := []byte(os.Getenv("POSTGRES_PASSWORD"))
		dboption := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, password)
		db, err := sql.Open("postgres", dboption)
		if err != nil {
			responseError(serverError.NewFatalServerError(err), w)
			return
		}

		// トランザクションオープン
		tx, err := db.Begin()
		if err != nil {
			responseError(serverError.NewFatalServerError(err), w)
			return
		}
		defer tx.Rollback()

		// トランザクションをコールバックに渡してあげる
		serr := f(w, r, tx)
		if serr != nil {
			responseError(serr, w)
			return
		}

		// トランザクションコミット
		err = tx.Commit()
		if err != nil {
			responseError(serr, w)
			return
		}

		// トランザクションクローズ
		defer db.Close()
	})
}

// 共通処理の前にJWTトークンをチェックするラッパー
func wrapWithAuth(f wrappedFunc) http.Handler {
	return JwtMiddleware.Handler(wrap(f))
}
