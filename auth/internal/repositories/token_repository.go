package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/todo-app/internal/domain"
)

type TokenRepositoryInterface interface {
	New(userId string, ttl time.Duration, scope string) (*domain.Token, error)
	Insert(token *domain.Token) error
}

// Check that the plaintext token has been provided and is exactly 26 bytes long.
// func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
//     v.Check(tokenPlaintext != "", "token", "must be provided")
//     v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
// }

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

// New is a shortcut which creates a new Token struct and then inserts the
// data in the tokens table.
func (r *TokenRepository) New(userId string, ttl time.Duration, scope string) (*domain.Token, error) {
	token, err := domain.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)
	return token, err
}

// Insert adds the data for a specific token to the tokens table
func (r *TokenRepository) Insert(token *domain.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser deletes all tokens for a specific user and scope
func (r *TokenRepository) DeleteAllForUser(scope, userId string) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, scope, userId)
	return err
}
