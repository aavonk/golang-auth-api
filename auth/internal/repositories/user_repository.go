package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/todo-app/internal/domain"
)

type UserRepositoryInterface interface {
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
	GetById(id string) (*domain.User, error)
	GetForToken(tokenScope, tokenPlaintext string) (*domain.User, error)
	Update(user *domain.User) error
}

type UserRepo struct {
	db *sqlx.DB
}

// Responsible for mapping struct fields to the database table columns
type UserDBModel struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Activated bool      `db:"activated"`
	CreatedAt time.Time `db:"created_at"`
}

// Returns a domain user object, insuring that we interact with the domain object,
// and keep certain things that may be database specific out of the application code.
func (m *UserDBModel) ToDomain() *domain.User {
	return &domain.User{
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
func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE email=$1`
	user := UserDBModel{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Activated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return user.ToDomain(), nil

}

// GetById searches for and returns a user given an ID.
// If not found, or an error occurs, and empty user is returned.
// This can be checked by using the isEmpty method on domain.User structs
//

func (r *UserRepo) GetById(id string) (*domain.User, error) {
	query := `SELECT * FROM USERS WHERE id = $1`
	user := UserDBModel{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Activated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return user.ToDomain(), nil
}

// Create inserts a user to the database. It takes a domain.User object
// as it's only parameter, and returns it if the insert is successful
// otherwise, it returns an error
func (r *UserRepo) Create(user *domain.User) (*domain.User, error) {

	model := UserDBModel{}
	query := `INSERT INTO users (id, first_name, last_name, email, password, activated)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, first_name, last_name, email, password, activated, created_at`

	args := []interface{}{user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&model.ID,
		&model.FirstName,
		&model.LastName,
		&model.Email,
		&model.Password,
		&model.Activated,
		&model.CreatedAt,
	)

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return model.ToDomain(), nil
}

// GetForToken will retrieve the details of the user associated with a particular
// activation token.
func (r *UserRepo) GetForToken(tokenScope, tokenPlaintext string) (*domain.User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client
	// This is a byte *array* with length of 32, not a slice

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT users.id, users.created_at, users.first_name, users.last_name, users.email, users.password, users.activated
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user UserDBModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Activated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user.ToDomain(), nil
}

// Update will update the details for a specific user. It will also check
// for a violation of the `users_email_key` constraint when preforming the
// update.
func (r *UserRepo) Update(user *domain.User) error {
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, email = $3, password = $4, activated = $5
	WHERE id = $6
	RETURNING id, first_name, last_name, email, password, activated, created_at`

	args := []interface{}{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Activated,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Activated,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}

	}
	return nil
}
