package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"UserBookApi/internal/application"
	"UserBookApi/internal/exception"

	"github.com/google/uuid"
)

var errBoom = errors.New("boom")

func TestRegister_Success(t *testing.T) {
	repo := newFakeRepo()
	uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})

	out, err := uc.Execute(context.Background(), application.RegisterUserInput{
		Name: "Joao", Email: "j@x.com", Password: "123456",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if out.UserID == uuid.Nil {
		t.Fatal("UserID veio vazio")
	}
	if repo.saveCalls != 1 {
		t.Fatalf("Save chamado %d vezes, esperado 1", repo.saveCalls)
	}

	saved := repo.users["j@x.com"]
	if saved.ID() != out.UserID {
		t.Error("o id retornado é diferente do id salvo")
	}
	if saved.PasswordHash() == "123456" {
		t.Error("a senha foi salva em texto puro")
	}
}

func TestRegister_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		in   application.RegisterUserInput
		want error
	}{
		{"nome vazio", application.RegisterUserInput{Name: "   ", Email: "j@x.com", Password: "123456"}, exception.ErrEmptyName},
		{"email vazio", application.RegisterUserInput{Name: "Joao", Email: "  ", Password: "123456"}, exception.ErrEmailEmpty},
		{"senha vazia", application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: ""}, exception.ErrEmptyPassword},
		{"senha so com espacos", application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: "      "}, exception.ErrEmptyPassword},
		{"senha curta", application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: "123"}, exception.ErrPasswordToSmall},
		{"senha longa demais", application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: strings.Repeat("a", 73)}, exception.ErrPasswordTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})

			_, err := uc.Execute(context.Background(), tt.in)

			if !errors.Is(err, tt.want) {
				t.Fatalf("erro = %v, esperado %v", err, tt.want)
			}
			if repo.saveCalls != 0 {
				t.Error("não deveria salvar com entrada inválida")
			}
		})
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := newFakeRepo()
	uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})
	in := application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: "123456"}

	if _, err := uc.Execute(context.Background(), in); err != nil {
		t.Fatalf("primeiro registro falhou: %v", err)
	}

	_, err := uc.Execute(context.Background(), in)

	if !errors.Is(err, exception.ErrEmailAlreadyExists) {
		t.Fatalf("erro = %v, esperado ErrEmailAlreadyExists", err)
	}
	if repo.saveCalls != 1 {
		t.Errorf("Save chamado %d vezes, esperado 1", repo.saveCalls)
	}
}

func TestRegister_NormalizesEmail(t *testing.T) {
	repo := newFakeRepo()
	uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})

	_, err := uc.Execute(context.Background(), application.RegisterUserInput{
		Name: "Joao", Email: "  J@X.com ", Password: "123456",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if _, ok := repo.users["j@x.com"]; !ok {
		t.Error("o email deveria ser salvo em minúsculas e sem espaços")
	}
}

func TestRegister_DependencyErrors(t *testing.T) {
	in := application.RegisterUserInput{Name: "Joao", Email: "j@x.com", Password: "123456"}

	t.Run("erro no hasher", func(t *testing.T) {
		repo := newFakeRepo()
		uc := application.NewRegisterUserUseCase(repo, &fakeHasher{hashErr: errBoom})

		_, err := uc.Execute(context.Background(), in)

		if !errors.Is(err, errBoom) {
			t.Fatalf("erro = %v, esperado errBoom", err)
		}
		if repo.saveCalls != 0 {
			t.Error("não deveria salvar se o hash falhou")
		}
	})

	t.Run("erro ao verificar email", func(t *testing.T) {
		repo := newFakeRepo()
		repo.existsErr = errBoom
		uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})

		_, err := uc.Execute(context.Background(), in)

		if !errors.Is(err, errBoom) {
			t.Fatalf("erro = %v, esperado errBoom", err)
		}
		if repo.saveCalls != 0 {
			t.Error("não deveria salvar se a verificação falhou")
		}
	})

	t.Run("erro ao salvar", func(t *testing.T) {
		repo := newFakeRepo()
		repo.saveErr = errBoom
		uc := application.NewRegisterUserUseCase(repo, &fakeHasher{})

		_, err := uc.Execute(context.Background(), in)

		if !errors.Is(err, errBoom) {
			t.Fatalf("erro = %v, esperado errBoom", err)
		}
	})
}