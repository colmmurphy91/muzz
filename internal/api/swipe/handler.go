package swipe

import (
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/api/response"
	"github.com/colmmurphy91/muzz/internal/entity"
	swipeService "github.com/colmmurphy91/muzz/internal/usecase/swipe"
)

type Handler struct {
	logger       *zap.SugaredLogger
	swipeService *swipeService.Service
}

func NewHandler(logger *zap.SugaredLogger, swipeService *swipeService.Service) *Handler {
	return &Handler{logger: logger, swipeService: swipeService}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/swipe", h.swipe)
}

func (h *Handler) swipe(w http.ResponseWriter, r *http.Request) {
	var swipe entity.Swipe

	if err := json.NewDecoder(r.Body).Decode(&swipe); err != nil {
		response.RenderErrorResponse(w, "failed", err)
		return
	}

	if err := swipe.Validate(); err != nil {
		response.RenderErrorResponse(w, "failed to validate", err)
		return
	}

	userID := r.Context().Value("user").(int)
	if userID != swipe.UserID {
		response.RenderErrorResponse(w, "forbidden", entity.ErrForbidden)
		return
	}

	matchResponse, err := h.swipeService.Swipe(r.Context(), swipe.UserID, swipe.TargetID, swipe.Preference)
	if err != nil {
		response.RenderErrorResponse(w, "failed to save swipe", err)
		return
	}

	response.RenderResponse(w, matchResponse, http.StatusCreated)
}
