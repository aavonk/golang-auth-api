package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) HashPassword() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	u.Password = string(pass)
	return nil
}

// Prepare generates a unique uuid and trims the space off the name and email
// fields of the user object
func (u *User) Prepare() {
	u.ID = uuid.New()
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

}

// ToHTTPResponse insures that the users password is never included in the response
func (u *User) ToHTTPResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Activated: u.Activated,
		CreatedAt: u.CreatedAt,
	}
}
