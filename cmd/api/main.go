package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"UserBookApi/internal/application"
	"UserBookApi/internal/infrastructure/postgres"
	"UserBookApi/internal/infrastructure/security"
	"UserBookApi/internal/interfaces/httpapi"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	keyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if keyPath == "" {
		log.Fatal("JWT_PRIVATE_KEY_PATH is required")
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	issuer, err := security.NewJWTIssuer(keyPEM, time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewUserRepository(db)
	hasher := security.BcryptHasher{}

	register := application.NewRegisterUserUseCase(repo, hasher)
	login := application.NewLoginUseCase(repo, hasher, issuer)
	handler := httpapi.NewHandler(register, login)

	r := chi.NewRouter()
	r.Post("/users", handler.Register)
	r.Post("/login", handler.Login)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}