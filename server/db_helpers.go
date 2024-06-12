package server

import (
	"database/sql"
	"fmt"

	"github.com/tscolari/servicetools/database"
)

// openDB will open a database connection based on the given
// configuration.
func openDB(config *database.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ToConnectStr())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	if config.ConnMaxLifeTime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifeTime)
	}

	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}

	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}
