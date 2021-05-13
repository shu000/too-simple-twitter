package user

import (
	"testing"
)

// 作成したユーザの情報が与えたパラメータと一致すること
func TestAuthorization(t *testing.T) {
	userId := "iamuser"
	name := "Mr. User"
	password := "pass1111"

	newUser := CreateUser(userId, name, password)

	if newUser.UserId != userId {
		t.Errorf("UserId不一致")
	}

	if newUser.Name != name {
		t.Errorf("Name不一致")
	}

	if !newUser.Authorization(password) {
		t.Errorf("PasswordHash不一致")
	}
}
