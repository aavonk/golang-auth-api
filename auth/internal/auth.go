package internal

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type JWTClaims struct {
	UserId uuid.UUID
	Email  string
}

func NewToken(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserId,
		"email":   claims.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
