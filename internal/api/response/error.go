package response

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/colmmurphy91/muzz/internal/api/model"
	"github.com/colmmurphy91/muzz/internal/entity"
	"github.com/colmmurphy91/muzz/internal/usecase/auth"
)

func RenderErrorResponse(w http.ResponseWriter, msg string, err error) {
	resp := model.ErrorResponse{Error: msg}
	status := http.StatusInternalServerError

	var validationErrs validation.Errors

	switch {
	case errors.As(err, &validationErrs):
		status = http.StatusBadRequest
		resp.Error = "Validation failed"
		resp.Validations = validationErrs
	case errors.Is(err, entity.ErrUserNotFound):
		status = http.StatusNotFound
		resp.Reason = "User does not exist"
	case errors.Is(err, entity.ErrEmailAlreadyExists), errors.Is(err, entity.ErrMatchAlreadyExists):
		status = http.StatusConflict
		resp.Reason = "already exists"
	case errors.Is(err, auth.ErrPasswordDoesNotMatch):
		status = http.StatusUnauthorized
		resp.Reason = "password does not match"
	case errors.Is(err, entity.ErrForbidden):
		status = http.StatusForbidden
		resp.Reason = msg
	}

	RenderResponse(w, resp, status)
}
