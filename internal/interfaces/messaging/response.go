package messaging

import (
	"errors"
	"log"

	"UserBookApi/internal/exception"
)

type Response struct {
	OK    bool       `json:"ok"`
	Data  any        `json:"data,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func fail(code, msg string) Response {
	return Response{Error: &ErrorBody{Code: code, Message: msg}}
}

func failFrom(err error) Response {
	switch {
	case errors.Is(err, exception.ErrEmailAlreadyExists):
		return fail("EMAIL_ALREADY_EXISTS", err.Error())
	case errors.Is(err, exception.ErrInvalidCredentials):
		return fail("INVALID_CREDENTIALS", err.Error())
	case errors.Is(err, exception.ErrEmptyName),
		errors.Is(err, exception.ErrEmailEmpty),
		errors.Is(err, exception.ErrEmptyPassword),
		errors.Is(err, exception.ErrPasswordToSmall),
		errors.Is(err, exception.ErrPasswordTooLong):
		return fail("INVALID_INPUT", err.Error())
	default:
		log.Printf("internal error: %v", err)
		return fail("INTERNAL", "internal error")
	}
}