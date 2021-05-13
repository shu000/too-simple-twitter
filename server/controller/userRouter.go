package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"too-simple-twitter/server/model/user"
	"too-simple-twitter/server/repository/userRepository"
	"too-simple-twitter/server/util/serverError"

	"github.com/gorilla/mux"
)

// ユーザリソースのAPIを定義
func UseUserRouter(router *mux.Router) {
	/*
		ユーザ新規作成
		Method: POST
		Path: /user
		Body: { user_id: string, name: string, password: string }
		ResponseCode: 201
		ResponseBody: User
	*/
	router.Handle("/user", wrap(createUser)).Methods(http.MethodPost)

	/*
		自分自身（Authorizationヘッダのトークンで認証されたユーザ）のプロフィールを取得
		Method: GET
		Path: /user/me
		Header: Authorhization
		ResponseCode: 200
		ResponseBody: UserProfile
	*/
	router.Handle("/user/me", wrapWithAuth(getMe)).Methods(http.MethodGet)

	/*
		ユーザをuser_idで前方一致検索
		Method: GET
		Path: /user
		Query: { user_id: string }
		Header: Authorhization
		ResponseCode: 200
		ResponseBody: []User
	*/
	router.Handle("/user", wrapWithAuth(getUsersById)).Methods(http.MethodGet)

	/*
		ユーザをフォローする
		Method: PATCH
		Path: /user/follow
		Body: { target_user_id: string }
		Header: Authorhization
		ResponseCode: 204
		ResponseBody: nil
	*/
	router.Handle("/user/follow", wrapWithAuth(follow)).Methods(http.MethodPatch)
}

/*
  ユーザ新規作成のリクエストボディ
*/
type createUserBody struct {
	UserId   string `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

/*
  ユーザ新規作成。検証用。
*/
func createUser(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	var body createUserBody
	err := parseRequestBody(&body, w, r)
	if err != nil {
		return err
	}

	if len(body.UserId) <= 0 || len(body.UserId) > 80 {
		return serverError.NewServerError(401, "ユーザIDは80文字までです")
	}

	if len(body.Name) <= 0 || len(body.Name) > 80 {
		return serverError.NewServerError(401, "ユーザ名は80文字までです")
	}

	if body.Password == "" {
		return serverError.NewServerError(401, "パスワードを入力してください")
	}

	newUser := user.CreateUser(body.UserId, body.Name, body.Password)
	err = userRepository.Save(tx, newUser)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)

	return nil
}

/*
	自分自身のプロフィールを取得
*/
func getMe(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	loginUserId, err := GetUserIdFromJWT(r)
	if err != nil {
		return err
	}

	targetUserProfile, err := userRepository.GetOneAsUserProfile(tx, loginUserId)
	if err != nil {
		return err
	}

	if targetUserProfile == nil {
		return serverError.NewServerError(409, "対象のユーザが見つかりません")
	}

	json.NewEncoder(w).Encode(targetUserProfile)

	return nil
}

/*
  ユーザをIDでLIKE検索
*/
func getUsersById(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	query := r.URL.Query()
	id := query.Get("user_id")

	users, err := userRepository.GetMany(tx, id)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(users)
	return nil
}

/*
  ユーザフォローのリクエストボディ
*/
type followBody struct {
	TargetUserId string `json:"target_user_id"`
}

/*
  ユーザをフォローする
*/
func follow(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	var body followBody
	err := parseRequestBody(&body, w, r)
	if err != nil {
		return err
	}

	loginUserId, err := GetUserIdFromJWT(r)
	if err != nil {
		return err
	}

	loginUserProfile, err := userRepository.GetOneAsUserProfile(tx, loginUserId)
	if err != nil {
		return err
	}

	if loginUserProfile == nil {
		return serverError.NewServerError(409, "ユーザ情報の取得に失敗しました")
	}

	targetUserProfile, err := userRepository.GetOneAsUserProfile(tx, body.TargetUserId)
	if err != nil {
		return err
	}

	if loginUserProfile == nil {
		return serverError.NewServerError(409, "フォロー対象者が見つかりません")
	}

	err = user.Follow(loginUserProfile, targetUserProfile)
	if err != nil {
		return err
	}

	err = userRepository.SaveFollowingsAndFollowed(tx, loginUserProfile)
	if err != nil {
		return err
	}

	// フォローされた側のSaveは呼ばない。テーブル構造上不要なため。
	// （Routerがテーブルの実装に依存していることになるので気持ち悪い）
	// err = userRepository.SaveFollowingsAndFollowed(tx, targetUserProfile)
	// if err != nil {
	// 	return err
	// }

	w.WriteHeader(http.StatusNoContent)

	return nil
}
