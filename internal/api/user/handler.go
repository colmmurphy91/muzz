package user

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/api/response"
	"github.com/colmmurphy91/muzz/internal/usecase/user"
)

type Handler struct {
	logger      *zap.SugaredLogger
	userManager *user.Manager
}

func NewHandler(logger *zap.SugaredLogger, userM *user.Manager) *Handler {
	return &Handler{logger: logger, userManager: userM}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/user/create", h.createUser)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	createUser, err := h.userManager.CreateUser(r.Context())
	if err != nil {
		response.RenderErrorResponse(w, "error creating user", err)
	}

	response.RenderResponse(w, createUser, http.StatusCreated)
}
