package application_test

import (
	"context"
	"errors"
	"testing"

	"UserBookApi/internal/application"
	"UserBookApi/internal/domain"
	"UserBookApi/internal/exception"
)

// seedUser coloca um usuário direto no repositório fake, já com o hash do fakeHasher.
func seedUser(t *testing.T, repo *fakeRepo, email, plainPassword string) domain.User {
	t.Helper()

	u, err := domain.NewUser("Joao", email, hashPrefix+plainPassword)
	if err != nil {
		t.Fatalf("falha ao criar usuário de teste: %v", err)
	}
	repo.users[u.Email()] = u
	return u
}

func TestLogin_Success(t *testing.T) {
	repo := newFakeRepo()
	user := seedUser(t, repo, "j@x.com", "123456")
	uc := application.NewLoginUseCase(repo, &fakeHasher{}, &fakeIssuer{})

	out, err := uc.Execute(context.Background(), application.LoginInput{
		Email: "j@x.com", Password: "123456",
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	want := "token-" + user.ID().String()
	if out.Token != want {
		t.Errorf("token = %q, esperado %q", out.Token, want)
	}
}

func TestLogin_EmailIsCaseInsensitive(t *testing.T) {
	repo := newFakeRepo()
	seedUser(t, repo, "j@x.com", "123456")
	uc := application.NewLoginUseCase(repo, &fakeHasher{}, &fakeIssuer{})

	_, err := uc.Execute(context.Background(), application.LoginInput{
		Email: "  J@X.COM ", Password: "123456",
	})

	if err != nil {
		t.Fatalf("login deveria funcionar com outra caixa: %v", err)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	tests := []struct {
		name string
		in   application.LoginInput
	}{
		{"senha errada", application.LoginInput{Email: "j@x.com", Password: "errada"}},
		{"email inexistente", application.LoginInput{Email: "outro@x.com", Password: "123456"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			seedUser(t, repo, "j@x.com", "123456")
			uc := application.NewLoginUseCase(repo, &fakeHasher{}, &fakeIssuer{})

			out, err := uc.Execute(context.Background(), tt.in)

			// os dois casos devem dar o MESMO erro, para não revelar quais emails existem
			if !errors.Is(err, exception.ErrInvalidCredentials) {
				t.Fatalf("erro = %v, esperado ErrInvalidCredentials", err)
			}
			if out.Token != "" {
				t.Error("não deveria emitir token")
			}
		})
	}
}

func TestLogin_DependencyErrors(t *testing.T) {
	t.Run("erro no repositório não vira credencial inválida", func(t *testing.T) {
		repo := newFakeRepo()
		repo.findErr = errBoom
		uc := application.NewLoginUseCase(repo, &fakeHasher{}, &fakeIssuer{})

		_, err := uc.Execute(context.Background(), application.LoginInput{
			Email: "j@x.com", Password: "123456",
		})

		if !errors.Is(err, errBoom) {
			t.Fatalf("erro = %v, esperado errBoom", err)
		}
		if errors.Is(err, exception.ErrInvalidCredentials) {
			t.Error("falha de infraestrutura não pode ser tratada como senha errada")
		}
	})

	t.Run("erro ao emitir token", func(t *testing.T) {
		repo := newFakeRepo()
		seedUser(t, repo, "j@x.com", "123456")
		uc := application.NewLoginUseCase(repo, &fakeHasher{}, &fakeIssuer{err: errBoom})

		_, err := uc.Execute(context.Background(), application.LoginInput{
			Email: "j@x.com", Password: "123456",
		})

		if !errors.Is(err, errBoom) {
			t.Fatalf("erro = %v, esperado errBoom", err)
		}
	})
}