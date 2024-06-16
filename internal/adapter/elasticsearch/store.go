package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/colmmurphy91/muzz/internal/entity"
)

type User struct {
	client *esv7.Client
	index  string
}

func NewUser(client *esv7.Client) *User {
	return &User{
		client: client,
		index:  "users",
	}
}

type indexedUser struct {
	ID       int             `json:"id"`
	Email    string          `json:"email"`
	Name     string          `json:"name"`
	Gender   string          `json:"gender"`
	Age      int             `json:"age"`
	Location entity.Location `json:"location"`
}

func (u *User) Index(ctx context.Context, user entity.User) error {
	body := indexedUser{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Gender:   user.Gender,
		Age:      user.Age,
		Location: user.Location,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return fmt.Errorf("failed to encode body: %w", err)
	}

	req := esv7api.IndexRequest{
		Index:      u.index,
		DocumentID: fmt.Sprint(body.ID),
		Body:       &buf,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, u.client)
	if err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to index %w", err)
	}

	io.Copy(io.Discard, resp.Body) //nolint: errcheck

	return nil
}

func (u *User) SearchOthers(ctx context.Context, params entity.SearchParams) ([]entity.User, error) {
	boolQuery := map[string]interface{}{
		"must_not": []map[string]interface{}{},
		"filter":   []map[string]interface{}{},
	}

	for _, userID := range params.ExcludeUserIDs {
		boolQuery["must_not"] = append(boolQuery["must_not"].([]map[string]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"id": userID,
			},
		})
	}

	if params.MinAge.Valid || params.MaxAge.Valid {
		ageRange := map[string]interface{}{}
		if params.MinAge.Valid {
			ageRange["gte"] = params.MinAge.Int64
		}
		if params.MaxAge.Valid {
			ageRange["lte"] = params.MaxAge.Int64
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"age": ageRange,
			},
		})
	}

	if params.Gender.Valid {
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"gender": params.Gender.String,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	resp, err := u.client.Search(
		u.client.Search.WithContext(ctx),
		u.client.Search.WithIndex(u.index),
		u.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("search query failed: %s", resp.String())
	}

	var r struct {
		Hits struct {
			Hits []struct {
				Source indexedUser `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	users := make([]entity.User, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		users[i] = entity.User{
			ID:       hit.Source.ID,
			Email:    hit.Source.Email,
			Name:     hit.Source.Name,
			Gender:   hit.Source.Gender,
			Age:      hit.Source.Age,
			Location: hit.Source.Location,
		}
	}

	return users, nil
}
