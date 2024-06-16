package match

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
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

func (s *Store) CreateMatch(ctx context.Context, match entity.Match) (entity.Match, error) {
	query := `
		INSERT INTO matches (user1_id, user2_id, match_id)
		VALUES (:user1_id, :user2_id, :match_id)
	`

	result, err := s.db.NamedExecContext(ctx, query, match)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.Match{}, fmt.Errorf("match already exists: %w", entity.ErrMatchAlreadyExists)
		}

		return entity.Match{}, fmt.Errorf("failed to create match: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.Match{}, fmt.Errorf("failed to extract id: %w", err)
	}

	match.ID = int(id)

	return match, nil
}
