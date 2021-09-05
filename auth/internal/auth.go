package internal

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func NewSession(r *http.Request, sessionName string) (*sessions.Session, error) {
	sess, err := SessionStore.Get(r, sessionName)

	if err != nil {
		return nil, err
	}

	sess.Options = &sessions.Options{
		// Path:     "/",
		MaxAge:   86400 * 7, // one week
		Secure:   true,
		HttpOnly: true,
	}

	return sess, nil
}
