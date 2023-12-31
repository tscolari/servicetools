package server

import (
	"context"
	"fmt"

	"github.com/tscolari/servicetools/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewWithDB returns a WithDB object configured with the given config.
func NewWithDB(databaseConfig *database.Config) (*WithDB, error) {
	db, err := gorm.Open(postgres.Open(databaseConfig.ToConnectStr()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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
	BaseDB *gorm.DB
}

// DB returns an usable database connection.
func (s *WithDB) DB(ctx context.Context) *gorm.DB {
	if s.BaseDB == nil {
		return nil
	}

	return s.BaseDB.Session(&gorm.Session{})
}

// ConfigureDatabase is the hook used by the cmd package to inject the
// WithDB object. Services using WithDB must overwrite this method.
func (s *WithDB) ConfigureDatabase(*WithDB) {
	panic("ConfigureDatabase must be implemented")
}
