package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)
	hasher := security.BcryptHasher{}
	register := application.NewRegisterUserUseCase(repo, hasher)
	handler := httpapi.NewHandler(register)

	r := chi.NewRouter()
	r.Post("/users", handler.Register)

	log.Fatal(http.ListenAndServe(":8080", r))
}