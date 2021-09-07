package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/todo-app/internal/domain"
	"github.com/todo-app/pkg/logger"
)

type UserRepositoryInterface interface {
	GetByEmail(email string) domain.User
	Create(user *domain.User) (domain.User, error)
	GetById(id string) domain.User
}

type UserRepo struct {
	db *sqlx.DB
}

// Responsible for mapping struct fields to the database table columns
type UserDBModel struct {
	ID             uuid.UUID `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	EmailConfirmed bool      `db:"email_confirmed"`
}

// Returns a domain user object, insuring that we interact with the domain object,
// and keep certain things that may be database specific out of the application code.
func (m *UserDBModel) ToDomain() domain.User {
	return domain.User{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Password:  m.Password,
	}
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// Get by email queries the DB for a certain user by email,
// returning the user if found, and returning and empty user
// domain model if not found.
//
// @example:
//  usr := repository.GetByEmail("email@email.com")
// 	if usr.IsEmpty() {
// 		... "handle empty/error case"
//	}
func (r *UserRepo) GetByEmail(email string) domain.User {
	user := UserDBModel{}

	// Get returns an error if the result set is empty
	err := r.db.Get(&user, "SELECT * FROM users WHERE email=$1", email)

	if err != nil {
		// No records were found, return empty user object
		return domain.User{}
	}

	return user.ToDomain()
}

// GetById searches for and returns a user given an ID.
// If not found, or an error occurs, and empty user is returned.
// This can be checked by using the isEmpty method on domain.User structs
//
// @example:
//  usr := repository.GetById("someId")
// 	if usr.IsEmpty() {
// 		... "handle empty/error case"
//	}
func (r *UserRepo) GetById(id string) domain.User {
	user := UserDBModel{}

	err := r.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)

	if err != nil {
		logger.Error.Printf("Failed GetUserById Query. Reason: %v", err)
		return domain.User{}
	}

	return user.ToDomain()
}

// Create inserts a user to the database. It takes a domain.User object
// as it's only parameter, and returns it if the insert is successful
// otherwise, it returns an error
//
// @Note:
// Before Saving it to the database, it hashes the password. If the hashing is done
// elsewhere, the password will then be double hashed and unable to check if a user has
// given a valid password
func (r *UserRepo) Create(user *domain.User) (domain.User, error) {
	user.HashPassword()

	model := &UserDBModel{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Password:       user.Password,
		EmailConfirmed: false,
	}

	_, err := r.db.NamedExec(`INSERT INTO users (id, first_name, last_name, email, password, email_confirmed)
	 VALUES (:id, :first_name, :last_name, :email, :password, :email_confirmed)`, model)

	if err != nil {
		return domain.User{}, err
	}

	// Find out how to return the user
	return model.ToDomain(), nil
}
