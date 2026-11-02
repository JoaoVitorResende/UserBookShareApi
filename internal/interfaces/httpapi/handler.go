package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"UserBookApi/internal/application"
	"UserBookApi/internal/exception"
)

type Handler struct {
	register *application.RegisterUserUseCase
}

func NewHandler(register *application.RegisterUserUseCase) *Handler {
	return &Handler{register: register}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	out, err := h.register.Execute(r.Context(), application.RegisterUserInput{
		Name: req.Name, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": out.UserID.String()})
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exception.ErrEmailAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, exception.ErrEmptyName),
		errors.Is(err, exception.ErrEmailEmpty),
		errors.Is(err, exception.ErrEmptyPassword),
		errors.Is(err, exception.ErrPasswordToSmall),
		errors.Is(err, exception.ErrPasswordTooLong):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}