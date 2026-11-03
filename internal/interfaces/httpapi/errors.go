package httpapi

import (
	"errors"
	"log"
	"net/http"

	"UserBookApi/internal/exception"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exception.ErrEmailAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, exception.ErrInvalidCredentials):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, exception.ErrEmptyName),
		errors.Is(err, exception.ErrEmailEmpty),
		errors.Is(err, exception.ErrEmptyPassword),
		errors.Is(err, exception.ErrPasswordToSmall),
		errors.Is(err, exception.ErrPasswordTooLong):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}