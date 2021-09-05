package services

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/repositories"
)

type IdentityServiceInterface interface {
	HandleLogin(req *identity.LoginRequest) (domain.User, error)
	HandleRegister(potentialUser *domain.User) (domain.User, error)
}

type IdentityService struct {
	UserRepo repositories.UserRepositoryInterface
}

func NewIdentityService(db *sqlx.DB) *IdentityService {
	return &IdentityService{
		UserRepo: repositories.NewUserRepository(db),
	}
}

// Handle login will return a token if login criteria is met, otherwise it will return an
// error.
func (s *IdentityService) HandleLogin(req *identity.LoginRequest) (domain.User, error) {
	// Does a user with this email exist? If not, respond with error
	existingUser := s.UserRepo.GetByEmail(req.Email)

	// if user is emoty, then no user was found -- invalid credentials
	if existingUser.IsEmpty() {
		return domain.User{}, errors.New("invalid credentials")
	}

	// Compare the passwords of the stored user & the supplied password
	err := identity.ComparePasswords([]byte(existingUser.Password), []byte(req.Passsword))

	if err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}

	// If passwords are same, we're good.

	if err != nil {
		return domain.User{}, errors.New("error generating token")
	}

	return existingUser, nil
}

func (s *IdentityService) HandleRegister(potentialUser *domain.User) (domain.User, error) {
	// Search for existing user
	found := s.UserRepo.GetByEmail(potentialUser.Email)

	if !found.IsEmpty() {
		return domain.User{}, errors.New("account already exists")
	}
	var err error
	potentialUser.Prepare()
	err = potentialUser.Validate()

	if err != nil {
		return domain.User{}, err
	}

	newUser, err := s.UserRepo.Create(potentialUser)

	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}
