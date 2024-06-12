package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tscolari/servicetools/database"
)

// NewWithRDB returns a WithRDB object configured with the given config.
func NewWithRDB(config *database.Config) (*WithRDB, error) {
	db, err := openDB(config)
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	return &WithRDB{
		BaseRDB: db,
	}, nil
}

// WithRDB defines a struct that has capability to access the database.
// Although this has no semantic restriction, the idea of this is to provide
// a second database access (on top of the WithDB) that is meant for read-only operations.
type WithRDB struct {
	BaseRDB *sql.DB
}

// ConfigureReaderDatabase is the hook used by the cmd package to inject the
// WithRDB object in the host struct. This must be implemented by the host struct.
func (s *WithRDB) ConfigureReaderDatabase(*WithRDB) {
	panic("ConfigureReaderDatabase must be implemented")
}

// RDB returns an usable DB connection, meant for read-only operations.
func (s *WithRDB) RDB(ctx context.Context) *sql.DB {
	if s.BaseRDB == nil {
		return nil
	}

	return s.BaseRDB
}
