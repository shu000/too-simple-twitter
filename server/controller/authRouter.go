package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"too-simple-twitter/server/repository/userRepository"
	"too-simple-twitter/server/util/serverError"

	"github.com/gorilla/mux"
)

// ログイン/ログアウトのAPIを定義
func UseAuthRouter(router *mux.Router) {
	/*
		ログイン
		Method: POST
		Path: /auth/login
		Body: { user_id: string, password: string }
		ResponseCode: 200
		ResponseBody: User
	*/
	router.Handle("/auth/login", wrap(login)).Methods(http.MethodPost)
}

// ログインのリクエストボディ
type loginRequestBody struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

// ログインのレスポンスボディ
type loginResponseBody struct {
	Token string `json:"token"`
}

/*
  ログイン
*/
func login(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	var body loginRequestBody
	err := parseRequestBody(&body, w, r)
	if err != nil {
		return err
	}

	if body.UserId == "" || body.Password == "" {
		return serverError.NewUnauthorizedServerError()
	}

	targetUser, err := userRepository.GetOne(tx, body.UserId)
	if err != nil {
		return err
	}

	if targetUser == nil {
		return serverError.NewUnauthorizedServerError()
	}

	if !targetUser.Authorization(body.Password) {
		return serverError.NewUnauthorizedServerError()
	}

	w.WriteHeader(http.StatusOK)

	res := new(loginResponseBody)
	res.Token = getJWT(targetUser.UserId)

	json.NewEncoder(w).Encode(res)

	return nil
}
