package messaging

import (
	"UserBookApi/internal/exception"
	"errors"
	"net/http"
)

type Reply struct {
	Status int    `json:"status"`
	Data   any    `json:"data"`
	Error  string `json:"error"`
}

func statusFromError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, exception.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, exception.ErrEmailAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}