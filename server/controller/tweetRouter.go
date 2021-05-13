package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"too-simple-twitter/server/model/tweet"
	"too-simple-twitter/server/repository/tweetRepository"
	"too-simple-twitter/server/repository/userRepository"
	"too-simple-twitter/server/util/serverError"

	"github.com/gorilla/mux"
)

// ツイートリソースのAPIを定義
func UseTweetRouter(router *mux.Router) {
	/*
		ツイートする
		Method: POST
		Path: /tweet
		Body: { contents: string }
		Header: Authorization
		ResponseCode: 201
		ResponseBody: Tweet
	*/
	router.Handle("/tweet", wrapWithAuth(postTweet)).Methods(http.MethodPost)

	/*
		タイムラインを取得する
		Method: GET
		Path: /tweet
		Header: Authorization
		ResponseCode: 200
		ResponseBody: Tweets
	*/
	router.Handle("/tweet", wrapWithAuth(getTimeline)).Methods(http.MethodGet)

}

/*
  ツイートのリクエストボディ
*/
type postTweetRequestBody struct {
	Contents string `json:"contents"`
}

/*
	ツイート
*/
func postTweet(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	var body postTweetRequestBody
	err := parseRequestBody(&body, w, r)
	if err != nil {
		return err
	}

	tweetUserId, err := GetUserIdFromJWT(r)
	if err != nil {
		return err
	}

	tweetUser, err := userRepository.GetOne(tx, tweetUserId)
	if err != nil {
		return err
	}

	if tweetUser == nil {
		return serverError.NewServerError(409, "対象のユーザが見つかりません")
	}

	if body.Contents == "" || len(body.Contents) > 140 {
		return serverError.NewServerError(401, "ツイートは140文字以下までです")
	}

	newTweet, err := tweet.CreateTweet(tweetUser, body.Contents)
	if err != nil {
		return err
	}

	err = tweetRepository.Save(tx, newTweet)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTweet)

	return nil
}

/*
  タイムライン取得
	TODO: ページネーション
*/
func getTimeline(w http.ResponseWriter, r *http.Request, tx *sql.Tx) serverError.ServerError {
	loginUserId, err := GetUserIdFromJWT(r)
	if err != nil {
		return err
	}

	loginUser, err := userRepository.GetOne(tx, loginUserId)
	if err != nil {
		return err
	}

	if loginUser == nil {
		return serverError.NewServerError(409, "対象のユーザが見つかりません")
	}

	tweets, err := tweetRepository.GetMany(tx, loginUser)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(tweets)
	return nil
}
