package user

import (
	"testing"
)

/*
  フォロー済みの場合エラーを返すこと
*/
func TestFollowThrowsErrorWhenAlreadyFollowing(t *testing.T) {
	following := NewUserProfile("following", "following", "pass", Users{}, Users{})
	followed := NewUserProfile("followed", "followed", "pass", Users{}, Users{})

	following.Follow(followed.GetAsUser())
	followed.FollowedBy(following.GetAsUser())

	err := Follow(following, followed)

	if err == nil {
		t.Errorf("フォロー済みエラー未検出")
	}
}

/*
  フォロー中が０人のとき成功すること
*/
func TestFollowSuccessWhenFollowingIsEmpty(t *testing.T) {
	following := NewUserProfile("following", "following", "pass", Users{}, Users{})
	followed := NewUserProfile("followed", "followed", "pass", Users{}, Users{})

	err := Follow(following, followed)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(following.Following) != 1 || !following.Following[0].EqualsTo(followed.GetAsUser()) {
		t.Errorf("フォロー中が増えていない")
	}

	if len(following.Following) != 1 || !followed.Followed[0].EqualsTo(following.GetAsUser()) {
		t.Errorf("フォロワーが増えていない")
	}
}

/*
  フォロー中が１人以上のとき成功すること
*/
func TestFollowSuccessWhenFollowingOthers(t *testing.T) {
	following := NewUserProfile("following", "following", "pass", Users{}, Users{})
	followed := NewUserProfile("followed", "followed", "pass", Users{}, Users{})

	user1 := NewUserProfile("user1", "user1", "pass", Users{}, Users{})

	err := Follow(following, user1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = Follow(user1, following)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = Follow(followed, user1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = Follow(user1, followed)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = Follow(following, followed)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(following.Following) != 2 || !following.Following[1].EqualsTo(followed.GetAsUser()) {
		t.Errorf("フォロー中が増えていない")
	}

	if len(following.Following) != 2 || !followed.Followed[1].EqualsTo(following.GetAsUser()) {
		t.Errorf("フォロワーが増えていない")
	}
}
