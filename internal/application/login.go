package application

import (
	"context"
	"errors"

	"UserBookApi/internal/domain"
	"UserBookApi/internal/exception"

	"github.com/google/uuid"
)

type TokenIssuer interface {
	Issue(userID uuid.UUID) (string, error)
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type LoginUseCase struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenIssuer
}

func Login(users UserRepository, hasher PasswordHasher, tokens TokenIssuer) *LoginUseCase {
	return &LoginUseCase{users: users, hasher: hasher, tokens: tokens}
}

func (uc *LoginUseCase) Execute(ctx context.Context, in LoginInput) (LoginOutput, error) {
	user, err := uc.users.FindByEmail(ctx, domain.NormalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, exception.ErrUserNotFound) {
			return LoginOutput{}, exception.ErrInvalidCredentials
		}
		return LoginOutput{}, err
	}

	if !uc.hasher.Compare(user.PasswordHash(), in.Password) {
		return LoginOutput{}, exception.ErrInvalidCredentials
	}

	token, err := uc.tokens.Issue(user.ID())
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{Token: token}, nil
}
