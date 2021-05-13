package controller

import (
	"too-simple-twitter/server/util/serverError"
	"net/http"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

// JWTがAuthorizationヘッダに存在するかチェック
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// JWTトークンからユーザIDを取得
func GetUserIdFromJWT(r *http.Request) (string, serverError.ServerError) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", serverError.NewUnauthorizedServerError()
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["user_id"].(string)
	if ok {
		return userId, nil
	}

	return "", serverError.NewUnauthorizedServerError()
}

// JWTを生成する
func getJWT(userId string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["user_id"] = userId
	// tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}
