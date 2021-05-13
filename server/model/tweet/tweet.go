/*
  ツイートパッケージ
*/
package tweet

import "time"

/*
  ツイート
*/
type Tweet struct {
	TweetId       string    `db:"tweet_id" json:"tweet_id"`
	TweetUserId   string    `db:"tweet_user_id" json:"tweet_user_id"`
	TweetedTime   time.Time `db:"tweeted_time" json:"tweeted_time"`
	Contents      string    `db:"contents" json:"contents"`
	IsDeleted     bool      `db:"id_deleted" json:"id_deleted"`
	TweetUserName string    `json:"tweet_user_name"`
}

// ユーザのスライス
type Tweets []*Tweet

// Tweetコンストラクタ
func NewTweet(tweetId string, tweetUserId string, tweetedTime time.Time, contents string, isDeleted bool, tweetUserName string) *Tweet {
	tweet := new(Tweet)
	tweet.TweetId = tweetId
	tweet.TweetUserId = tweetUserId
	tweet.TweetedTime = tweetedTime
	tweet.Contents = contents
	tweet.IsDeleted = isDeleted
	tweet.TweetUserName = tweetUserName

	return tweet
}

// IDが一致したらtrue
func (tweet *Tweet) EqualsTo(target *Tweet) bool {
	return tweet.TweetId == target.TweetId
}

// ツイートを削除
func (tweet *Tweet) Delete() {
	tweet.IsDeleted = false
}
