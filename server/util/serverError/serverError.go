package serverError

import (
	"fmt"
	"os"
)

// 独自エラーインターフェース
type ServerError interface {
	error
	Code() int
}

// HTTPエラーコードを含むerrorsインターフェース実装
type structServerError struct {
	code    int    // HTTPエラーコード
	message string // エラーメッセージ
}

// errorsインターフェース実装
func (err *structServerError) Error() string {
	return err.message
}

// HTTPエラーコード
func (err *structServerError) Code() int {
	return err.code
}

// ServerErrorコンストラクタ
func NewServerError(code int, message string) *structServerError {
	err := new(structServerError)
	err.code = code
	err.message = message

	return err
}

/*
	500エラー用のコンストラクタ。予期してないエラーをキャッチしたときに使う。
  コンソールに元のエラーメッセージを吐いたのち、ユーザ向けのエラーメッセージを返す
  TODO: ちゃんとしたloggerつかう
*/
func NewFatalServerError(err error) *structServerError {
	fmt.Fprintln(os.Stderr, err)
	return NewServerError(500, "サーバーエラーです。しばらく経ってからやり直してください。")
}

// 401エラー用のコンストラクタ。メッセージが共通なので。
func NewUnauthorizedServerError() *structServerError {
	return NewServerError(401, "認証エラーです。ログインしてください。")
}
