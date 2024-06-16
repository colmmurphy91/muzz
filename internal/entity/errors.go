package entity

import (
	"errors"
)

var ErrMatchAlreadyExists = errors.New("match already exists")

var (
	ErrEmailAlreadyExists = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user does not exists")
)

var ErrForbidden = errors.New("forbidden")

var ErrInvalidParam = errors.New("invalid param")
