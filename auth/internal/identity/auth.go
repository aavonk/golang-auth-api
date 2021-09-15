package identity

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	IdentitySessionName   = "user-session"
	SessionStore          *sessions.CookieStore
	cookies               *securecookie.SecureCookie
)

func init() {
	SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // one week
		Secure:   true,
		HttpOnly: true,
	}

	cookies = securecookie.New([]byte(os.Getenv("SESSION_KEY")), nil)

}

type LoginRequest struct {
	Email     string `json:"email"`
	Passsword string `json:"password"`
}

type JWTClaims struct {
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.StandardClaims
}

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func ComparePasswords(hashedPassword, suppliedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, suppliedPassword)
}

// TODO: Add Expiration date on token

func newToken(claims JWTClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ExtractClaimsFromToken parses a JWT token into a JWTClaims struct
// which should include a UserId and Email value. If there is an error
// parsing the token, the error is return and an empty JWTClaims struct is returned
// as well.
func ExtractClaimsFromToken(tokenString string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Return Claims object
		return *claims, nil

	} else {
		// return an error
		return JWTClaims{}, err
	}
}

// TODO: Add Expiration date on cookie
func SetCookie(w http.ResponseWriter, user *domain.User) error {
	token, err := newToken(JWTClaims{
		UserId: user.ID,
		Email:  user.Email,
	})

	if err != nil {
		return err
	}
	cookieValue := map[string]string{
		"token": token,
	}

	encoded, err := cookies.Encode("auth-session", cookieValue)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:  "auth-session",
		Value: encoded,
		Path:  "/",
		// Domain:     "",
		// Expires:    time.Time{},
		// MaxAge:     0,
		Secure:   true,
		HttpOnly: true,
		// SameSite:   0,
		Unparsed: []string{},
	}
	http.SetCookie(w, cookie)

	return nil
}

// GetTokenFromCookie extracts a cookie from the request named "auth-session"
// If not present, it will return an error and an emtpy string.
// Once we verify that the cookie is present, we decode it, which should
// contain a key value pair of "jwt": "{JSON web token}".
func GetTokenFromCookie(r *http.Request) (string, error) {

	value := make(map[string]string)

	cookie, err := r.Cookie("auth-session")
	if err != nil {
		logger.Error.Println("Error getting cookie", err)
		return "", err
	}

	if err := cookies.Decode("auth-session", cookie.Value, &value); err != nil {
		logger.Error.Println("Error Decoding cookie", err)
		return "", err
	}
	token := value["token"]
	return token, nil
}

var UserCtxKey = &authContextKey{"user"}

type authContextKey struct {
	name string
}

func GetClaimsFromContext(ctx context.Context) (JWTClaims, bool) {
	raw, ok := ctx.Value(UserCtxKey).(JWTClaims)
	return raw, ok
}
