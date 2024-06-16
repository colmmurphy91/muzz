// nolint
package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/colmmurphy91/muzz/internal/adapter/mysql/user/model"
	"github.com/colmmurphy91/muzz/internal/usecase/auth/mocks"
)

func TestAuthService_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserFetcher := mocks.NewMockuserFetcher(ctrl)
	secretKey := "testsecretkey"
	authService := NewAuthService(secretKey, mockUserFetcher)

	tests := []struct {
		name          string
		email         string
		password      string
		setupMock     func()
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "successful authentication",
			email:    "test@example.com",
			password: "password123",
			setupMock: func() {
				mockUserFetcher.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(model.User{
						Email:    "test@example.com",
						Password: "password123",
					}, nil)
			},
			expectedToken: "", // We will verify the token format instead of exact match
			expectedErr:   nil,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "password123",
			setupMock: func() {
				mockUserFetcher.EXPECT().
					FindByEmail(gomock.Any(), "notfound@example.com").
					Return(model.User{}, fmt.Errorf("user not found"))
			},
			expectedToken: "",
			expectedErr:   fmt.Errorf("failed to retrieve user: user not found"),
		},
		{
			name:     "incorrect password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				mockUserFetcher.EXPECT().
					FindByEmail(gomock.Any(), "test@example.com").
					Return(model.User{
						Email:    "test@example.com",
						Password: "password123",
					}, nil)
			},
			expectedToken: "",
			expectedErr:   ErrPasswordDoesNotMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			token, err := authService.Authenticate(context.Background(), tt.email, tt.password)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				// Verify the token format instead of exact match
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(secretKey), nil
				})
				assert.NoError(t, err)
				assert.True(t, parsedToken.Valid)
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				assert.True(t, ok)
				assert.Equal(t, tt.email, claims["email"])
				assert.WithinDuration(t, time.Unix(int64(claims["exp"].(float64)), 0), time.Now().Add(time.Hour*72), time.Minute)
			}
		})
	}
}
