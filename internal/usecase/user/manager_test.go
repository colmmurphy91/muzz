package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/colmmurphy91/muzz/internal/adapter/mysql/user/model"
	"github.com/colmmurphy91/muzz/internal/entity"
	"github.com/colmmurphy91/muzz/internal/usecase/user/mocks"
)

func TestManager_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserCreator := mocks.NewMockuserCreator(ctrl)
	mockUserIndexer := mocks.NewMockuserIndexer(ctrl)

	manager := NewManager(mockUserCreator, mockUserIndexer)

	tests := []struct {
		name          string
		setupMocks    func()
		expectedUser  entity.User
		expectedError error
	}{
		{
			name: "successful creation and indexing",
			setupMocks: func() {
				mockUserCreator.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{
						ID:       1,
						Email:    "test-email@muzz.com",
						Password: "password123",
						Name:     "Test User",
						Gender:   "Male",
						Age:      30,
					}, nil)

				mockUserIndexer.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedUser: entity.User{
				ID:       1,
				Email:    "test-email@muzz.com",
				Password: "password123",
				Name:     "Test User",
				Gender:   "Male",
				Age:      30,
			},
			expectedError: nil,
		},
		{
			name: "user creation failure",
			setupMocks: func() {
				mockUserCreator.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{}, errors.New("creation failed"))
			},
			expectedUser:  entity.User{},
			expectedError: errors.New("failed to create user: creation failed"),
		},
		{
			name: "indexing failure",
			setupMocks: func() {
				mockUserCreator.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(model.User{
						ID:       1,
						Email:    "test-email@muzz.com",
						Password: "password123",
						Name:     "Test User",
						Gender:   "Male",
						Age:      30,
					}, nil)

				mockUserIndexer.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(errors.New("indexing failed"))
			},
			expectedUser:  entity.User{},
			expectedError: errors.New("failed to index: indexing failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			user, err := manager.CreateUser(context.Background())

			assert.Equal(t, tt.expectedUser, user)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
