package tweetRepository

import (
	"database/sql"
	"time"

	"too-simple-twitter/server/model/tweet"
	"too-simple-twitter/server/model/user"
	"too-simple-twitter/server/util/serverError"
)

/*
  ユーザをINSERTする
  NOTE: 今回は要件にないが、本来UPSERTにすべき
*/
func Save(tx *sql.Tx, newTweet *tweet.Tweet) serverError.ServerError {
	_, err := tx.Exec("INSERT INTO tweets (tweet_id, tweet_user_id, tweeted_time, contents, is_deleted) VALUES ($1, $2, $3, $4, $5)",
		newTweet.TweetId,
		newTweet.TweetUserId,
		newTweet.TweetedTime,
		newTweet.Contents,
		newTweet.IsDeleted,
	)

	if err != nil {
		return serverError.NewFatalServerError(err)
	}

	return nil
}

func GetMany(tx *sql.Tx, targetUser *user.User) (tweet.Tweets, serverError.ServerError) {
	rows, err := tx.Query(`
		SELECT DISTINCT ON (a.tweet_id, a.tweeted_time) a.tweet_id, a.tweet_user_id, a.tweeted_time, a.contents, a.is_deleted, c.name  FROM tweets a
		LEFT JOIN followers b
		ON b.followed_user_id = a.tweet_user_id
		LEFT JOIN users c
		ON a.tweet_user_id = c.user_id
		WHERE a.is_deleted = false AND (a.tweet_user_id = $1 OR b.following_user_id = $2)
		ORDER BY a.tweeted_time DESC
	`, targetUser.UserId, targetUser.UserId)

	if err != nil {
		return tweet.Tweets{}, serverError.NewFatalServerError(err)

	}
	defer rows.Close()

	tweets := tweet.Tweets{}
	var tweetId, tweetUserId, contents, tweetUserName string
	var tweetedTime time.Time
	var isDeleted bool
	for rows.Next() {
		err := rows.Scan(&tweetId, &tweetUserId, &tweetedTime, &contents, &isDeleted, &tweetUserName)
		if err != nil {
			return tweet.Tweets{}, serverError.NewFatalServerError(err)

		}

		tweets = append(tweets, tweet.NewTweet(tweetId, tweetUserId, tweetedTime, contents, isDeleted, tweetUserName))
	}

	err = rows.Err()
	if err != nil {
		return tweet.Tweets{}, serverError.NewFatalServerError(err)
	}

	return tweets, nil
}
