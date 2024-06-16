package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
	Reason      string            `json:"reason,omitempty"`
}
