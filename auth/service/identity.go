package service

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/todo-app/internal"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
)

type LoginRequest struct {
	Email     string `json:"email"`
	Passsword string `json:"password"`
}

type IdentityServiceInterface interface {
	HandleLogin(req *LoginRequest) (domain.User, error)
}

type IdentityService struct {
	UserRepo domain.UserRepository
}

func NewIdentityService(db *sqlx.DB) *IdentityService {
	return &IdentityService{
		UserRepo: internal.NewUserRepository(db),
	}
}

// Handle login will return a token if login criteria is met, otherwise it will return an
// error.
func (s *IdentityService) HandleLogin(req *LoginRequest) (string, error) {
	// Does a user with this email exist? If not, respond with error
	existingUser := s.UserRepo.GetByEmail(req.Email)

	// if user is emoty, then no user was found -- invalid credentials
	if existingUser.IsEmpty() {
		return "", errors.New("invalid credentials")
	}

	// Compare the passwords of the stored user & the supplied password
	err := identity.ComparePasswords([]byte(existingUser.Password), []byte(req.Passsword))

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// If passwords are same, we're good. Send them a JWT in a cookie

	token, err := identity.NewToken(identity.JWTClaims{
		UserId: existingUser.ID,
		Email:  existingUser.Email,
	})

	if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil
}