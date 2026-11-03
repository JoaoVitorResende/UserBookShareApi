package httpapi

import (
	"encoding/json"
	"net/http"

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

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	out, err := h.useCase.Execute(r.Context(), application.RegisterUserInput{
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