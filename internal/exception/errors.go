package exception

import "errors"

var (
	ErrEmptyName          = errors.New("name can't be empty")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrPasswordToSmall    = errors.New("password to small")
	ErrEmailEmpty         = errors.New("email empty")
	ErrPasswordTooLong    = errors.New("password to large")
	ErrEmailAlreadyExists = errors.New("email already exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)
