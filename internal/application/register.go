package application

import (
	"context"

	"UserBookApi/internal/domain"
	"UserBookApi/internal/exception"

	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user domain.User) error
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserOutput struct {
	UserID uuid.UUID
}

type RegisterUserUseCase struct {
	users  UserRepository
	hasher PasswordHasher
}

func NewRegisterUserUseCase(users UserRepository, hasher PasswordHasher) *RegisterUserUseCase {
	return &RegisterUserUseCase{users: users, hasher: hasher}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, request RegisterUserInput) (RegisterUserOutput, error) {
	if err := domain.ValidatePlainPassword(request.Password); err != nil {
		return RegisterUserOutput{}, err
	}

	hash, err := uc.hasher.Hash(request.Password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	user, err := domain.NewUser(request.Name, request.Email, hash)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	exists, err := uc.users.ExistsByEmail(ctx, user.Email())
	if err != nil {
		return RegisterUserOutput{}, err
	}
	if exists {
		return RegisterUserOutput{}, exception.ErrEmailAlreadyExists
	}

	if err := uc.users.Save(ctx, user); err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{UserID: user.ID()}, nil
}
