package repositories

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrEditConflict   = errors.New("conflict submitting edit operation")
)
