package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lr *LoginRequest) Validate() error {
	return validation.ValidateStruct(
		lr,
		validation.Field(&lr.Email, validation.Required, is.Email),
		validation.Field(&lr.Password, validation.Required, validation.Length(6, 100)),
	)
}

type TokenResponse struct {
	Token string `json:"token"`
}
