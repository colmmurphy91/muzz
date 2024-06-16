package tools

import (
	"context"
	"fmt"
	"strings"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/colmmurphy91/muzz/tools/envvar"
)

// NewElasticSearch instantiates the ElasticSearch client using configuration defined in environment variables.
func NewElasticSearch(conf *envvar.Configuration) (es *esv7.Client, err error) {
	esHost := conf.Get("ES_HOST")
	esPort := conf.Get("ES_PORT")

	cfg := esv7.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", esHost, esPort),
		},
	}

	es, err = esv7.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster information: %w", err)
	}

	defer func() {
		err = res.Body.Close()
	}()

	return es, nil
}

func CreateUsersIndex(es *esv7.Client) error {
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"id": { "type": "integer" },
				"email": { "type": "keyword" },
				"name": { "type": "text" },
				"gender": { "type": "keyword" },
				"age": { "type": "integer" },
				"location": { "type": "geo_point" }
			}
		}
	}`

	req := esv7api.IndicesCreateRequest{
		Index: "users",
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}
