package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	UserId        string
	Token         string
	DateGenerated int64
	DateExpires   int64
}

func NewEmailVerificationItem(user *User) *EmailVerification {
	now := time.Now()

	return &EmailVerification{
		UserId:        user.ID.String(),
		Token:         uuid.New().String(),
		DateGenerated: now.Unix(),
		DateExpires:   now.AddDate(0, 0, 7).Unix(), // Expires one week from now
	}
}
