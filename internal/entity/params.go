package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	null "github.com/guregu/null/v5"
)

type SearchParams struct {
	ExcludeUserIDs []int
	MinAge         null.Int
	MaxAge         null.Int
	Gender         null.String
	Lat            float64
	Lon            float64
}

// IsZero determines whether the search arguments have values or not.
func (p SearchParams) IsZero() bool {
	return !p.MinAge.Valid && !p.MaxAge.Valid && !p.Gender.Valid
}

func (p SearchParams) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MinAge, validation.By(nonNegativeInt)),
		validation.Field(&p.MaxAge, validation.By(nonNegativeInt)),
		validation.Field(&p.Gender, validation.By(validGender)),
	)
}

// nolint:forcetypeassert
func nonNegativeInt(value interface{}) error {
	if value.(null.Int).Valid && value.(null.Int).Int64 < 0 {
		return validation.NewError("validation_non_negative", "must be non-negative")
	}

	return nil
}

// nolint:forcetypeassert
func validGender(value interface{}) error {
	gender := value.(null.String)
	if gender.Valid && gender.String != "" && gender.String != "male" && gender.String != "female" {
		return validation.NewError("validation_invalid_gender", "must be 'male', 'female', or empty")
	}

	return nil
}
