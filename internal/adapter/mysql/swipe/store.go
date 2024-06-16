package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/entity"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// GetUsersYesSwipes retrieves all swipes where the user swiped "YES"
func (s *Store) GetTargetsYesSwipes(ctx context.Context, userID int) (map[int]entity.Swipe, error) {
	swipes := []entity.Swipe{}
	query := `
		SELECT id, user_id, target_id, preference
		FROM swipes
		WHERE user_id = ? AND preference = 'YES'
	`

	err := s.db.SelectContext(ctx, &swipes, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find target swipes: %w", err)
	}

	swipeMap := make(map[int]entity.Swipe, len(swipes))
	for _, swipe := range swipes {
		swipeMap[swipe.TargetID] = swipe
	}

	return swipeMap, nil
}

// SaveSwipe saves a new swipe without overwriting existing records
func (s *Store) SaveSwipe(ctx context.Context, swipe entity.Swipe) error {
	query := `
		INSERT INTO swipes (user_id, target_id, preference)
		VALUES (:user_id, :target_id, :preference)
		ON DUPLICATE KEY UPDATE id = id
	`

	_, err := s.db.NamedExecContext(ctx, query, swipe)
	if err != nil {
		return fmt.Errorf("failed to save swipe: %w", err)
	}

	return nil
}

func (s *Store) GetUserSwipes(ctx context.Context, userID int) ([]entity.Swipe, error) {
	swipes := []entity.Swipe{}
	query := `
		SELECT id, user_id, target_id, preference
		FROM swipes
		WHERE user_id = ?
	`

	err := s.db.SelectContext(ctx, &swipes, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find target swipes: %w", err)
	}
	return swipes, nil
}
