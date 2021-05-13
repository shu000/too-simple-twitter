package user

/*
  ユーザ
  認証やユーザ名のリストとして利用
*/
type User struct {
	UserId       string `db:"user_id" json:"user_id"`
	Name         string `db:"name" json:"name"`
	PasswordHash string `db:"password_hash" json:"-"` // パスワードハッシュはクライアントには教えない
}

// ユーザのスライス
type Users []*User

// Userコンストラクタ
func NewUser(userId string, name string, passwordHash string) *User {
	user := new(User)
	user.UserId = userId
	user.Name = name
	user.PasswordHash = passwordHash

	return user
}

// IDが一致したらtrue
func (user *User) EqualsTo(target *User) bool {
	return user.UserId == target.UserId
}

// パスワードハッシュが一致したらtrue
func (user *User) Authorization(password string) bool {
	passwordHash := createHash(password)
	if user.PasswordHash == passwordHash {
		return true
	} else {
		return false
	}
}
