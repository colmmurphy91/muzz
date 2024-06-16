package service

import (
	"context"
	"fmt"

	"github.com/colmmurphy91/muzz/internal/entity"
)

//go:generate mockgen -source $GOFILE -destination mocks/mocks_${GOFILE} -package mocks

type swiper interface {
	GetTargetsYesSwipes(ctx context.Context, userID int) (map[int]entity.Swipe, error)
	SaveSwipe(ctx context.Context, swipe entity.Swipe) error
}

type matcher interface {
	CreateMatch(ctx context.Context, match entity.Match) (entity.Match, error)
}

type Service struct {
	swiper  swiper
	matcher matcher
}

func NewService(swipe swiper, match matcher) *Service {
	return &Service{
		swiper:  swipe,
		matcher: match,
	}
}

func (s *Service) Swipe(ctx context.Context, userID, target int, preference entity.Preference) (MatchResponse, error) {
	// Save the swipe
	err := s.swiper.SaveSwipe(ctx, entity.Swipe{
		UserID:     userID,
		TargetID:   target,
		Preference: preference,
	})
	if err != nil {
		return MatchResponse{}, fmt.Errorf("failed to save swipe: %w", err)
	}

	if len(preference) == len(entity.PreferenceNo) {
		return MatchResponse{Matched: false}, nil
	}

	// Check if it's a match
	yesSwipes, err := s.swiper.GetTargetsYesSwipes(ctx, target)
	if err != nil {
		return MatchResponse{}, fmt.Errorf("failed to get user's yes swipes: %w", err)
	}

	if _, aMatch := yesSwipes[userID]; aMatch {
		match := entity.Match{
			User1ID: userID,
			User2ID: target,
		}

		match.GenerateMatchID()

		createdMatch, err := s.matcher.CreateMatch(ctx, match)
		if err != nil {
			return MatchResponse{}, fmt.Errorf("failed to create match: %w", err)
		}

		return MatchResponse{
			Matched: true,
			MatchID: createdMatch.ID,
		}, nil
	}

	return MatchResponse{Matched: false}, nil
}

type MatchResponse struct {
	Matched bool `json:"matched"`
	MatchID int  `json:"matchID,omitempty"`
}
