package user

/*
  ユーザプロフィール
  ユーザにフォロー・フォロワーのスライスを付与したデータ

  ユーザ検索や認証時など、毎回フォロワースライスを保持するのは過剰なため、Userと分離した
  またUserProfileが存在しない場合ロジックがすべてRopositoryに集約してしまうのが嫌
*/
type UserProfile struct {
	User
	Following Users `json:"following"`
	Followed  Users `json:"followed"`
}

// ユーザプロフィールのスライス
type UserProfiles []*UserProfile

// UserProfileコンストラクタ
func NewUserProfile(userId string, name string, passwordHash string, following Users, followed Users) *UserProfile {
	userProfile := new(UserProfile)
	userProfile.UserId = userId
	userProfile.Name = name
	userProfile.PasswordHash = passwordHash
	userProfile.Following = following
	userProfile.Followed = followed

	return userProfile
}

// IDが一致したらtrue
func (userProfile *UserProfile) EqualsTo(target *UserProfile) bool {
	return userProfile.UserId == target.UserId
}

// UserProfileからUserを抽出
func (userProfile *UserProfile) GetAsUser() *User {
	user := new(User)
	user.UserId = userProfile.UserId
	user.Name = userProfile.Name
	user.PasswordHash = userProfile.PasswordHash

	return user
}

// ユーザをフォローする
func (userProfile *UserProfile) Follow(followed *User) error {
	userProfile.Following = append(userProfile.Following, followed)
	return nil
}

// ユーザにフォローされる
func (userProfile *UserProfile) FollowedBy(following *User) error {
	userProfile.Followed = append(userProfile.Followed, following)
	return nil
}

// 指定ユーザをフォローしていたらtrue
func (userProfile *UserProfile) IsFollowing(target *User) bool {
	for _, followed := range userProfile.Following {
		if followed.EqualsTo(target) {
			return true
		}
	}

	return false
}
