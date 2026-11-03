package messaging

import (
	"context"
	"encoding/json"

	"UserBookApi/internal/application"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterHandler struct {
	useCase *application.RegisterUserUseCase
}

func NewRegisterHandler(uc *application.RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{useCase: uc}
}

func (h *RegisterHandler) Handle(ctx context.Context, body []byte) Response {
	var req registerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return fail("INVALID_INPUT", "invalid body")
	}

	out, err := h.useCase.Execute(ctx, application.RegisterUserInput{
		Name: req.Name, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		return failFrom(err)
	}

	return Response{OK: true, Data: map[string]string{"id": out.UserID.String()}}
}