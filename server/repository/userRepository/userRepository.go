package userRepository

import (
	"database/sql"
	"fmt"
	"strings"

	"too-simple-twitter/server/model/user"
	"too-simple-twitter/server/util/serverError"
)

/*
  ユーザをINSERTする
  NOTE: 今回は要件にないが、本来UPSERTにすべき
*/
func Save(tx *sql.Tx, newUser *user.User) serverError.ServerError {
	_, err := tx.Exec("INSERT INTO users (user_id, name, password_hash) VALUES ($1, $2, $3)", newUser.UserId, newUser.Name, newUser.PasswordHash)
	if err != nil {
		return serverError.NewFatalServerError(err)
	}

	return nil
}

/*
  Userをidで取得する
*/
func GetOne(tx *sql.Tx, userId string) (*user.User, serverError.ServerError) {
	var got user.User
	err := tx.QueryRow("SELECT * FROM users WHERE user_id = $1", userId).Scan(&got.UserId, &got.Name, &got.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return new(user.User), serverError.NewFatalServerError(err)
	}

	return &got, nil
}

/*
  ユーザをidで前方一致検索する
*/
func GetMany(tx *sql.Tx, id string) (user.Users, serverError.ServerError) {
	rows, err := tx.Query("SELECT * FROM users WHERE user_id LIKE $1", id+"%")
	if err != nil {
		return user.Users{}, serverError.NewFatalServerError(err)

	}
	defer rows.Close()

	users := user.Users{}
	var userId, name, passwordHash string
	for rows.Next() {
		err := rows.Scan(&userId, &name, &passwordHash)
		if err != nil {
			return user.Users{}, serverError.NewFatalServerError(err)

		}

		users = append(users, user.NewUser(userId, name, passwordHash))
	}

	err = rows.Err()
	if err != nil {
		return user.Users{}, serverError.NewFatalServerError(err)

	}

	return users, nil
}

/*
  UserProfileをidで取得する
*/
func GetOneAsUserProfile(tx *sql.Tx, targetUserId string) (*user.UserProfile, serverError.ServerError) {
	// ユーザ取得
	var got user.UserProfile
	err := tx.QueryRow("SELECT * FROM users WHERE user_id = $1", targetUserId).Scan(&got.UserId, &got.Name, &got.PasswordHash)
	if err != nil {
		return new(user.UserProfile), serverError.NewFatalServerError(err)
	}

	// 未存在の場合はnilを返す
	if got.UserId == "" {
		return nil, nil
	}

	// フォロー中を取得
	followingUsers, serr := getFollowingUsers(tx, targetUserId)
	if err != nil {
		return new(user.UserProfile), serr
	}
	got.Following = followingUsers

	// フォロワーを取得
	followedUsers, serr := getFollowedUsers(tx, targetUserId)
	if err != nil {
		return new(user.UserProfile), serr
	}
	got.Followed = followedUsers

	return &got, nil
}

/*
  対象のフォロー中とフォロワーをすべて保存する
  全フォロー中とフォロワーをDELETE/INSERTしている
  引数のUserProfileの存在確認はしてない。RDBから取得済みのUserProfileを渡すこと。
*/
func SaveFollowingsAndFollowed(tx *sql.Tx, userProfile *user.UserProfile) serverError.ServerError {
	// すべてのフォロー中とフォロワーをDELETE
	_, err := tx.Exec("DELETE FROM followers WHERE following_user_id = $1 OR followed_user_id = $2", userProfile.UserId, userProfile.UserId)
	if err != nil {
		return serverError.NewFatalServerError(err)
	}

	// すべてのフォロー中とフォロワーをINSERTするためのvaluesを構築
	// UserProfile中のUserIdを渡している（パラメータを直で渡していない）と信じてSpringfで構築
	values := ""
	for _, following := range userProfile.Following {
		values = values + fmt.Sprintf("('%s', '%s'),", userProfile.UserId, following.UserId)
	}
	for _, followed := range userProfile.Followed {
		values = values + fmt.Sprintf("('%s', '%s'),", followed.UserId, userProfile.UserId)
	}
	values = strings.TrimRight(values, ",")

	// フォロー中もフォロワーも０人ならINSERTを実行せずに終了
	// TODO: 未テスト。フォロー解除実装時に検証すること。
	if values == "" {
		return nil
	}

	// すべてのフォロー中とフォロワーをINSERT
	_, err = tx.Exec("INSERT INTO followers (following_user_id, followed_user_id) VALUES " + values)
	if err != nil {
		return serverError.NewFatalServerError(err)
	}

	return nil
}

// フォロー中のユーザをすべて取得
func getFollowingUsers(tx *sql.Tx, targetUserId string) (user.Users, serverError.ServerError) {
	followingUsers := user.Users{}

	rows, err := tx.Query(`
    SELECT c.user_id, c.name, c.password_hash FROM users a
    RIGHT JOIN followers b
		ON b.following_user_id = a.user_id
		LEFT JOIN users c
		ON c.user_id = b.followed_user_id
		WHERE a.user_id = $1
	`, targetUserId)

	if err != nil {
		return user.Users{}, serverError.NewFatalServerError(err)
	}
	defer rows.Close()

	var userId, name, passwordHash string
	for rows.Next() {
		err := rows.Scan(&userId, &name, &passwordHash)
		if err != nil {
			return user.Users{}, serverError.NewFatalServerError(err)
		}

		followingUsers = append(followingUsers, user.NewUser(userId, name, passwordHash))
	}

	return followingUsers, nil
}

// フォローされているユーザをすべて取得
func getFollowedUsers(tx *sql.Tx, targetUserId string) (user.Users, serverError.ServerError) {
	followedUsers := user.Users{}

	rows, err := tx.Query(`
    SELECT c.user_id, c.name, c.password_hash FROM users a
    RIGHT JOIN followers b
		ON b.followed_user_id = a.user_id
		LEFT JOIN users c
		ON c.user_id = b.following_user_id
		WHERE a.user_id = $1
	`, targetUserId)
	if err != nil {
		return user.Users{}, serverError.NewFatalServerError(err)
	}
	defer rows.Close()

	var userId, name, passwordHash string
	for rows.Next() {
		err := rows.Scan(&userId, &name, &passwordHash)
		if err != nil {
			return user.Users{}, serverError.NewFatalServerError(err)
		}

		followedUsers = append(followedUsers, user.NewUser(userId, name, passwordHash))
	}

	return followedUsers, nil
}
