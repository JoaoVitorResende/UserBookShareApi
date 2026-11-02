package domain

import (
	"strings"

	"UserBookApi/internal/exception"

	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	name         string
	email        string
	passwordHash string
}

func ValidatePlainPassword(plain string) error {
	switch {
	case strings.TrimSpace(plain) == "":
		return exception.ErrEmptyPassword
	case len(plain) < 6:
		return exception.ErrPasswordToSmall
	case len(plain) > 72:
		return exception.ErrPasswordTooLong
	}
	return nil
}

func validateUserItens(name, email, passwordHash string) error{

    if name == "" {
		return exception.ErrEmptyName
	}
	if email == "" {
		return exception.ErrEmailEmpty
	}
	if passwordHash == "" {
		return exception.ErrEmptyPassword
	}

    return nil
}

func NewUser(name, email, passwordHash string) (User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
    err := validateUserItens(name, email, passwordHash)

    if err != nil{
        return User{}, err
    }

	return User{id: uuid.New(), name: name, email: email, passwordHash: passwordHash}, nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func RestoreUser(id uuid.UUID, name, email, passwordHash string) User {
	return User{id: id, name: name, email: email, passwordHash: passwordHash}
}

func (u User) ID() uuid.UUID        { return u.id }
func (u User) Name() string         { return u.name }
func (u User) Email() string        { return u.email }
func (u User) PasswordHash() string { return u.passwordHash }
