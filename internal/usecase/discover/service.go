package discover

import (
	"context"
	"fmt"

	"github.com/asmarques/geodist"

	"github.com/colmmurphy91/muzz/internal/entity"
)

type userDiscover interface {
	SearchOthers(ctx context.Context, params entity.SearchParams) ([]entity.User, error)
}

type swiper interface {
	GetUserSwipes(ctx context.Context, currentUserID int) ([]entity.Swipe, error)
}

type Service struct {
	userDiscover userDiscover
	swiper       swiper
}

func NewService(discover userDiscover, swipe swiper) *Service {
	return &Service{
		userDiscover: discover,
		swiper:       swipe,
	}
}

func (s *Service) DiscoverPeople(ctx context.Context, userID int, params entity.SearchParams) ([]entity.User, error) {
	swipes, err := s.swiper.GetUserSwipes(ctx, userID)
	if err != nil {
		return nil, err
	}

	var excludedIds = []int{userID}

	for _, swipe := range swipes {
		excludedIds = append(excludedIds, swipe.TargetID)
	}

	params.ExcludeUserIDs = excludedIds

	users, err := s.userDiscover.SearchOthers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search for others %w", err)
	}

	var userPoint = geodist.Point{Lat: params.Lat, Long: params.Lon}

	for i, user := range users {
		km := geodist.HaversineDistance(userPoint, geodist.Point{
			Lat:  user.Location.Lat,
			Long: user.Location.Lon,
		})
		users[i].DistanceFromMe = km
	}

	return users, nil
}
