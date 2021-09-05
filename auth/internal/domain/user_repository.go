package domain

type UserRepository interface {
	GetByEmail(email string) User
	Create(user *User) (User, error)
}
