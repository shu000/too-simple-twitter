package tweet

import (
	"time"

	"too-simple-twitter/server/model/user"
	"too-simple-twitter/server/util/serverError"
)

// 新規ツイートを作成
func CreateTweet(tweetUser *user.User, contents string) (*Tweet, serverError.ServerError) {
	uuid, err := GetUUID()
	if err != nil {
		return new(Tweet), serverError.NewFatalServerError(err)
	}

	return NewTweet(uuid, tweetUser.UserId, time.Now(), contents, false, tweetUser.Name), nil
}
