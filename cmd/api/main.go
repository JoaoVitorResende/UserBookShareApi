package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"UserBookApi/internal/application"
	"UserBookApi/internal/infrastructure/postgres"
	"UserBookApi/internal/infrastructure/security"
	"UserBookApi/internal/interfaces/httpapi"
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

	// infrastructure
	repo := postgres.NewUserRepository(db)
	hasher := security.BcryptHasher{}

	// application
	register := application.NewRegisterUserUseCase(repo, hasher)
	login := application.NewLoginUseCase(repo, hasher, issuer)

	// interfaces: HTTP (opcional, só sobe se HTTP_ADDR estiver definido)
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		srv := &http.Server{Addr: addr, Handler: httpapi.NewRouter(register, login)}
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("http: %v", err)
			}
		}()
		defer srv.Shutdown(context.Background())
	}

	// interfaces: Rabbit
	consumer := messaging.NewConsumer(ch)
	consumer.Register(messaging.QueueRegister, messaging.NewRegisterHandler(register))
	consumer.Register(messaging.QueueLogin, messaging.NewLoginHandler(login))

	if err := consumer.Run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("shutting down")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s é obrigatória", key)
	}
	return v
}