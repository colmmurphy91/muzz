package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/colmmurphy91/muzz/internal/entity"
	"github.com/colmmurphy91/muzz/internal/usecase/swipe/mocks"
)

func TestService_Swipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwiper := mocks.NewMockswiper(ctrl)
	mockMatcher := mocks.NewMockmatcher(ctrl)
	service := NewService(mockSwiper, mockMatcher)

	ctx := context.Background()
	userID := 1
	targetID := 2
	preferenceYes := entity.Preference("YES")
	preferenceNo := entity.Preference("NO")

	tests := []struct {
		name            string
		userID          int
		targetID        int
		preference      entity.Preference
		mockSwiper      func()
		mockMatcher     func()
		expectedResp    MatchResponse
		expectedError   error
		shouldCallMatch bool
	}{
		{
			name:       "successful swipe, no match",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceYes,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceYes,
				}).Return(nil)
				mockSwiper.EXPECT().GetUsersYesSwipes(ctx, gomock.Any()).Return(map[int]entity.Swipe{}, nil)
			},
			mockMatcher:     func() {},
			expectedResp:    MatchResponse{Matched: false},
			expectedError:   nil,
			shouldCallMatch: false,
		},
		{
			name:       "successful swipe, match found",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceYes,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceYes,
				}).Return(nil)
				mockSwiper.EXPECT().GetUsersYesSwipes(ctx, gomock.Any()).Return(map[int]entity.Swipe{
					userID: {UserID: targetID, TargetID: userID, Preference: preferenceYes},
				}, nil)
			},
			mockMatcher: func() {
				mockMatcher.
					EXPECT().
					CreateMatch(ctx, gomock.Any()).DoAndReturn(
					func(_ context.Context, match entity.Match) (entity.Match, error) {
						match.ID = 1
						return match, nil
					})
			},
			expectedResp:    MatchResponse{Matched: true, MatchID: 1},
			expectedError:   nil,
			shouldCallMatch: true,
		},
		{
			name:       "swipe no",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceNo,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceNo,
				}).Return(nil)
			},
			mockMatcher:     func() {},
			expectedResp:    MatchResponse{Matched: false},
			expectedError:   nil,
			shouldCallMatch: false,
		},
		{
			name:       "error saving swipe",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceYes,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceYes,
				}).Return(fmt.Errorf("db error"))
			},
			mockMatcher:     func() {},
			expectedResp:    MatchResponse{},
			expectedError:   fmt.Errorf("failed to save swipe: db error"),
			shouldCallMatch: false,
		},
		{
			name:       "error getting yes swipes",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceYes,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceYes,
				}).Return(nil)
				mockSwiper.EXPECT().GetUsersYesSwipes(ctx, gomock.Any()).Return(nil, fmt.Errorf("db error"))
			},
			mockMatcher:     func() {},
			expectedResp:    MatchResponse{},
			expectedError:   fmt.Errorf("failed to get user's yes swipes: db error"),
			shouldCallMatch: false,
		},
		{
			name:       "error creating match",
			userID:     userID,
			targetID:   targetID,
			preference: preferenceYes,
			mockSwiper: func() {
				mockSwiper.EXPECT().SaveSwipe(ctx, entity.Swipe{
					UserID:     userID,
					TargetID:   targetID,
					Preference: preferenceYes,
				}).Return(nil)
				mockSwiper.EXPECT().GetUsersYesSwipes(ctx, targetID).Return(map[int]entity.Swipe{
					userID: {UserID: userID, TargetID: targetID, Preference: preferenceYes},
				}, nil)
			},
			mockMatcher: func() {
				mockMatcher.EXPECT().CreateMatch(ctx, gomock.Any()).Return(entity.Match{}, fmt.Errorf("db error"))
			},
			expectedResp:    MatchResponse{},
			expectedError:   fmt.Errorf("failed to create match: db error"),
			shouldCallMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSwiper()

			if tt.shouldCallMatch {
				tt.mockMatcher()
			}

			resp, err := service.Swipe(ctx, tt.userID, tt.targetID, tt.preference)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
