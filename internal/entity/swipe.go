package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Preference string

type Swipe struct {
	ID         int        `json:"id" db:"id"`
	UserID     int        `json:"user_id" db:"user_id"`
	TargetID   int        `json:"target_id" db:"target_id"`
	Preference Preference `json:"preference" db:"preference"`
}

const (
	PreferenceYes Preference = "yes"
	PreferenceNo  Preference = "no"
)

func (s Swipe) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserID, validation.Required, validation.Min(1)),
		validation.Field(&s.TargetID, validation.Required, validation.Min(1)),
		validation.Field(&s.Preference, validation.Required, validation.In(PreferenceYes, PreferenceNo)),
	)
}
