package httpapi

import (
	"net/http"

	"UserBookApi/internal/application"

	"github.com/go-chi/chi/v5"
)

func NewRouter(register *application.RegisterUserUseCase, login *application.LoginUseCase) http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/users", NewRegisterHandler(register))
	r.Method(http.MethodPost, "/login", NewLoginHandler(login))
	return r
}