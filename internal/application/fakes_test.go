package application_test

import (
	"context"

	"UserBookApi/internal/domain"
	"UserBookApi/internal/exception"

	"github.com/google/uuid"
)

const hashPrefix = "hashed:"

// fakeRepo implementa application.UserRepository em memória.
type fakeRepo struct {
	users     map[string]domain.User // chave: e-mail
	existsErr error
	saveErr   error
	findErr   error
	saveCalls int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: make(map[string]domain.User)}
}

func (r *fakeRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if r.existsErr != nil {
		return false, r.existsErr
	}
	_, ok := r.users[email]
	return ok, nil
}

func (r *fakeRepo) Save(ctx context.Context, u domain.User) error {
	r.saveCalls++
	if r.saveErr != nil {
		return r.saveErr
	}
	r.users[u.Email()] = u
	return nil
}

func (r *fakeRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if r.findErr != nil {
		return domain.User{}, r.findErr
	}
	u, ok := r.users[email]
	if !ok {
		return domain.User{}, exception.ErrUserNotFound
	}
	return u, nil
}

// fakeHasher implementa application.PasswordHasher sem bcrypt (rápido e previsível).
type fakeHasher struct {
	hashErr error
}

func (h *fakeHasher) Hash(plain string) (string, error) {
	if h.hashErr != nil {
		return "", h.hashErr
	}
	return hashPrefix + plain, nil
}

func (h *fakeHasher) Compare(hash, plain string) bool {
	return hash == hashPrefix+plain
}

// fakeIssuer implementa application.TokenIssuer.
type fakeIssuer struct {
	err error
}

func (i *fakeIssuer) Issue(userID uuid.UUID) (string, error) {
	if i.err != nil {
		return "", i.err
	}
	return "token-" + userID.String(), nil
}