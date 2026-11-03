package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"UserBookApi/internal/application"
	"UserBookApi/internal/infrastructure/postgres"
	"UserBookApi/internal/infrastructure/security"
	"UserBookApi/internal/interfaces/messaging"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := mustEnv("DATABASE_URL")
	amqpURL := mustEnv("RABBITMQ_URL")
	keyPath := mustEnv("JWT_PRIVATE_KEY_PATH")

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	issuer, err := security.NewJWTIssuer(keyPEM, time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewUserRepository(db)
	hasher := security.BcryptHasher{}

	register := application.NewRegisterUserUseCase(repo, hasher)
	login := application.NewLoginUseCase(repo, hasher, issuer)

	consumer := messaging.NewConsumer(ch, register, login)
	if err := consumer.Run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("shutting down")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}