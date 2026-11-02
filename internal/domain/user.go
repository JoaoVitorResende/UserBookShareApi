package domain

import (
    "UserBookApi/internal/exception"
    "strings"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    userId       uuid.UUID
    name         string
    passwordHash string
    email        string
    bookQuantity int
}

func NewUser(name, password, email string) (User, error) {
    name = strings.TrimSpace(name)
    email = strings.TrimSpace(email)

    if err := validate(name, password, email); err != nil {
        return User{}, err
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return User{}, err
    }

    return User{
        userId:       uuid.New(),
        name:         name,
        email:        email,
        passwordHash: string(hash),
    }, nil
}

func (u User) CheckPassword(plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(plain)) == nil
}

func validate(name, password, email string) error {
    switch {
    case name == "":
        return exception.ErrEmptyName
    case email == "":
        return exception.ErrEmailEmpty
    case strings.TrimSpace(password) == "":
        return exception.ErrEmptyPassword
    case len(password) < 6:
        return exception.ErrPasswordToSmall
    case len(password) > 72: 
        return exception.ErrPasswordTooLong 
    }
    return nil
}