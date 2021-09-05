package identity

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/todo-app/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	IdentitySessionName = "user-session"
)

type LoginRequest struct {
	Email     string `json:"email"`
	Passsword string `json:"password"`
}

var SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type JWTClaims struct {
	UserId uuid.UUID
	Email  string
}

func newToken(claims JWTClaims) (string, error) {
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

func ComparePasswords(hashedPassword, suppliedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, suppliedPassword)
}

func newSession(r *http.Request, sessionName string) (*sessions.Session, error) {
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

func SetAndSaveSession(r *http.Request, w http.ResponseWriter, user domain.User) error {
	session, err := newSession(r, IdentitySessionName)

	if err != nil {
		return err
	}

	token, err := newToken(JWTClaims{
		UserId: user.ID,
		Email:  user.Email,
	})

	if err != nil {
		return err
	}

	// Save a JWT Key in the session cookie
	session.Values["jwt"] = token

	err = session.Save(r, w)

	if err != nil {
		return err
	}
	return nil
}
