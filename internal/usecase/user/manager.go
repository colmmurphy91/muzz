package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	faker "github.com/bxcodec/faker/v3"

	"github.com/colmmurphy91/muzz/internal/adapter/mysql/user/model"
	"github.com/colmmurphy91/muzz/internal/entity"
)

//go:generate mockgen -source $GOFILE -destination mocks/mocks_${GOFILE} -package mocks

type userCreator interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type userIndexer interface {
	Index(ctx context.Context, user entity.User) error
}

type Manager struct {
	userCreator userCreator
	userIndexer userIndexer
}

func NewManager(creator userCreator, indexer userIndexer) *Manager {
	return &Manager{userCreator: creator, userIndexer: indexer}
}

func (m *Manager) CreateUser(ctx context.Context) (entity.User, error) {
	user := generateFakeUser()

	dbUser, err := m.userCreator.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = dbUser.ID

	newUser := entity.User{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Name:     dbUser.Name,
		Password: dbUser.Password,
		Gender:   dbUser.Gender,
		Age:      dbUser.Age,
		Location: entity.Location{
			Lat: dbUser.Lat,
			Lon: dbUser.Lon,
		},
	}

	indexErr := m.userIndexer.Index(ctx, newUser)
	if indexErr != nil {
		return entity.User{}, fmt.Errorf("failed to index: %w", indexErr)
	}

	return newUser, nil
}

func generateFakeUser() model.User {
	p, _ := faker.RandomInt(18, 70, 1)

	name := faker.Name()
	sanitizedName := sanitizeName(name)
	email := fmt.Sprintf("%s@muzz.com", sanitizedName)

	return model.User{
		Email:    email,
		Password: "Password1",
		Name:     name,
		Gender:   faker.Gender(),
		Age:      p[0],
		Lon:      faker.Longitude(),
		Lat:      faker.Latitude(),
	}
}

func sanitizeName(name string) string {
	// Convert the name to lowercase
	name = strings.ToLower(name)
	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	// Remove invalid characters
	re := regexp.MustCompile(`[^a-z0-9-_]`)
	name = re.ReplaceAllString(name, "")

	return name
}
