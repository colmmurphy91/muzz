package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/colmmurphy91/muzz/internal/adapter/mysql/user/model"
)

//go:generate mockgen -source $GOFILE -destination mocks/mocks_${GOFILE} -package mocks

var ErrPasswordDoesNotMatch = errors.New("password does not match")

type userFetcher interface {
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

type Service struct {
	secretKey   string
	userFetcher userFetcher
}

func NewAuthService(secretKey string, fetcher userFetcher) *Service {
	return &Service{secretKey: secretKey, userFetcher: fetcher}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := s.userFetcher.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user.Password != password {
		return "", ErrPasswordDoesNotMatch
	}

	return s.generateJWT(email, user.ID)
}

func (s *Service) generateJWT(email string, id int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %w", err)
	}

	return signedString, nil
}
