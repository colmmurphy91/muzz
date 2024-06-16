package tools

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/colmmurphy91/muzz/tools/envvar"
)

// NewLocalMySQL creates a new MySQL database connection
func NewDBConnection(conf *envvar.Configuration) (*sqlx.DB, error) {
	databaseHost := conf.Get("DATABASE_HOST")
	databasePort := conf.Get("DATABASE_PORT")
	databaseUsername := conf.Get("DATABASE_USERNAME")
	databasePassword := conf.Get("DATABASE_PASSWORD")
	databaseName := conf.Get("DATABASE_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		databaseUsername,
		databasePassword,
		databaseHost,
		databasePort,
		databaseName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
