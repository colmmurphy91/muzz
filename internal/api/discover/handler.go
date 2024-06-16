package discover

import (
	"net/http"
	"strconv"

	"github.com/colmmurphy91/muzz/tools"

	chi "github.com/go-chi/chi/v5"
	null "github.com/guregu/null/v5"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/api/response"
	"github.com/colmmurphy91/muzz/internal/entity"
	"github.com/colmmurphy91/muzz/internal/usecase/discover"
)

type Handler struct {
	logger          *zap.SugaredLogger
	discoverService *discover.Service
}

func NewHandler(logger *zap.SugaredLogger, discoverService *discover.Service) *Handler {
	return &Handler{logger: logger, discoverService: discoverService}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/discover", h.discover)
}

// nolint:cyclop,forcetypeassert
func (h *Handler) discover(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(tools.CTXUserKey).(int)

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	params := entity.SearchParams{}

	if latStr == "" || lonStr == "" {
		http.Error(w, "lat and lon query parameters are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid value for lat", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid value for lon", http.StatusBadRequest)
		return
	}

	if minAgeParam := r.URL.Query().Get("min_age"); minAgeParam != "" {
		minAge, err := strconv.Atoi(minAgeParam)
		if err != nil {
			response.RenderErrorResponse(w, "invalid param", entity.ErrInvalidParam)

			return
		}

		params.MinAge = null.IntFrom(int64(minAge))
	}

	if maxAgeParam := r.URL.Query().Get("max_age"); maxAgeParam != "" {
		maxAge, err := strconv.Atoi(maxAgeParam)
		if err != nil {
			response.RenderErrorResponse(w, "invalid param", entity.ErrInvalidParam)

			return
		}

		params.MaxAge = null.IntFrom(int64(maxAge))
	}

	if genderParam := r.URL.Query().Get("gender"); genderParam != "" {
		params.Gender = null.StringFrom(genderParam)
	}

	if err := params.Validate(); err != nil {
		response.RenderErrorResponse(w, "failed to validate", err)
		return
	}

	params.Lat = lat
	params.Lon = lon

	people, err := h.discoverService.DiscoverPeople(r.Context(), userID, params)
	if err != nil {
		response.RenderErrorResponse(w, "failed to discover", err)
		return
	}

	response.RenderResponse(w, people, http.StatusOK)
}
