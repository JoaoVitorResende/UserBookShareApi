package postgres

import (
	"UserBookApi/internal/domain"
	"UserBookApi/internal/exception"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) Save(ctx context.Context, u domain.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`,
		u.ID(), u.Name(), u.Email(), u.PasswordHash())

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return exception.ErrEmailAlreadyExists
	}
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var (
		id           uuid.UUID
		name, mail   string
		passwordHash string
	)

	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password_hash FROM users WHERE email = $1`, email,
	).Scan(&id, &name, &mail, &passwordHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, exception.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return domain.RestoreUser(id, name, mail, passwordHash), nil
}
