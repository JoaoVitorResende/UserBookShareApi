package messaging

import (
	"context"
	"encoding/json"

	"UserBookApi/internal/application"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginHandler struct {
	useCase *application.LoginUseCase
}

func NewLoginHandler(uc *application.LoginUseCase) *LoginHandler {
	return &LoginHandler{useCase: uc}
}

func (h *LoginHandler) Handle(ctx context.Context, body []byte) Response {
	var req loginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return fail("INVALID_INPUT", "invalid body")
	}

	out, err := h.useCase.Execute(ctx, application.LoginInput{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		return failFrom(err)
	}

	return Response{OK: true, Data: map[string]string{"token": out.Token}}
}