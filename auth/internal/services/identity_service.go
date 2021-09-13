package services

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/todo-app/internal/domain"
	"github.com/todo-app/internal/identity"
	"github.com/todo-app/internal/repositories"
)

type IdentityServiceInterface interface {
	HandleLogin(req *identity.LoginRequest) (*domain.User, error)
	HandleRegister(potentialUser *domain.User) (*domain.User, error)
	GetUserById(id string) (*domain.User, error)
}

type IdentityService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewIdentityService(db *sqlx.DB) *IdentityService {
	return &IdentityService{
		userRepo: repositories.NewUserRepository(db),
	}
}

// Handle login will return a token if login criteria is met, otherwise it will return an
// error.
func (s *IdentityService) HandleLogin(req *identity.LoginRequest) (*domain.User, error) {
	// Does a user with this email exist? If not, respond with error
	existingUser, err := s.userRepo.GetByEmail(req.Email)

	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Compare the passwords of the stored user & the supplied password
	err = identity.ComparePasswords([]byte(existingUser.Password), []byte(req.Passsword))

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// If passwords are same, we're good.

	if err != nil {
		return nil, errors.New("error generating token")
	}

	return existingUser, nil
}

// HandleRegister processes the register request. If the request is denied - the email is taken,
// the password doesn't meet strength criteria, etc. - and empty user object and an error is returned.
// If it passes the checks, a new user is inserted into the database and returned.
// The potentialUser param must be a pointer to a domain.User struct so that the password
// can be hashed.
func (s *IdentityService) HandleRegister(potentialUser *domain.User) (*domain.User, error) {
	// There's no need to check if a user already exists with the email
	// address given. When the repository attempts to insert a user into the DB
	// it handles the error that would come from attempting to insert a user
	// with the same email, and returns an internal error(ErrDuplicateEmail).

	// Hash the password before sending to DB
	potentialUser.HashPassword()

	newUser, err := s.userRepo.Create(potentialUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *IdentityService) GetUserById(id string) (*domain.User, error) {

	return s.userRepo.GetById(id)

}
