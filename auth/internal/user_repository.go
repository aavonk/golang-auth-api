package internal

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/todo-app/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

type UserDBModel struct {
	ID             uuid.UUID `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	EmailConfirmed bool      `db:"email_confirmed"`
}

// TODO: Remove password from domain model so its not passed through code and accidentally returned?
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
