package user

import (
	"too-simple-twitter/server/util/serverError"
)

// パスワードのハッシュ値を生成しながらUserを作成する
func CreateUser(userId string, name string, password string) *User {
	user := new(User)
	user.UserId = userId
	user.Name = name
	user.PasswordHash = createHash(password)

	return user
}

// ユーザをフォローする
func Follow(following *UserProfile, followed *UserProfile) serverError.ServerError {
	followedUser := followed.GetAsUser()

	if following.IsFollowing(followedUser) {
		return serverError.NewServerError(409, "すでにフォローしています")
	}

	following.Follow(followedUser)
	followed.FollowedBy(following.GetAsUser())

	return nil
}
