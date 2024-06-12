package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tscolari/servicetools/database"
)

// NewWithDB returns a WithDB object configured with the given config.
func NewWithDB(config *database.Config) (*WithDB, error) {
	db, err := openDB(config)
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	return &WithDB{
		BaseDB: db,
	}, nil
}

// WithDB defines a struct that has capability to access the database.
// Access to the database object should be done using the `DB()` method.
type WithDB struct {
	// BaseDB contains a source connection to the database.
	// This is not safe to be used directly, and it's exposed only for
	// the purposes of making testing (and modifications) easier.
	BaseDB *sql.DB
}

// DB returns an usable database object.
// If no database configured, it will return nil.
func (s *WithDB) DB(ctx context.Context) *sql.DB {
	if s.BaseDB == nil {
		return nil
	}

	return s.BaseDB
}

// ConfigureDatabase is the hook used by the cmd package to inject the
// WithDB object in the host struct. This must be implemented by the host struct.
func (s *WithDB) ConfigureDatabase(*WithDB) {
	panic("ConfigureDatabase must be implemented")
}
