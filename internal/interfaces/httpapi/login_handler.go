package httpapi

import (
	"encoding/json"
	"net/http"

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

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	out, err := h.useCase.Execute(r.Context(), application.LoginInput{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": out.Token})
}