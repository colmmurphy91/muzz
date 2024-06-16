package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/colmmurphy91/muzz/internal/adapter/mysql/user/model"
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

func (s *Store) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	query := `INSERT INTO users(email, password, name, gender,age, lat, lon) 
	VALUES (:email, :password, :name, :gender, :age, :lat, :lon)`

	result, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return model.User{}, fmt.Errorf("email already exists: %w", entity.ErrEmailAlreadyExists)
		}

		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, fmt.Errorf("failed to retrieve id: %w", err)
	}

	user.ID = int(id)

	return user, nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := "SELECT id, email, password, name, gender, age, lat, lon FROM users WHERE email = ?"

	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, entity.ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}
