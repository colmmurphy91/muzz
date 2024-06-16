package login

import (
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/api/login/model"
	"github.com/colmmurphy91/muzz/internal/api/response"
	"github.com/colmmurphy91/muzz/internal/usecase/auth"
)

type Handler struct {
	logger      *zap.SugaredLogger
	authService *auth.Service
}

func NewHandler(logger *zap.SugaredLogger, auth *auth.Service) *Handler {
	return &Handler{logger: logger, authService: auth}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/login", h.login)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		response.RenderErrorResponse(w, "Invalid request payload", err)
		return
	}

	err = loginRequest.Validate()
	if err != nil {
		response.RenderErrorResponse(w, "Validation failed", err)
		return
	}

	token, err := h.authService.Authenticate(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		response.RenderErrorResponse(w, "failed to authenticate", err)
		return
	}

	response.RenderResponse(w, model.TokenResponse{
		Token: token,
	}, http.StatusOK)
}
