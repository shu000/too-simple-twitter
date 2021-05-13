package user

import (
	"crypto/sha512"
	"fmt"
)

// ユーザにまつわるヘルパー関数たち

// 与えられた文字列のSHA512のHEX表現を返す
func createHash(password string) string {
	hash := sha512.Sum512([]byte(password))
	hex := fmt.Sprintf("%x", hash)
	return hex
}
