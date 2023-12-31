package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tscolari/servicetools/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewWithRDBFromJSONConfig(jsonConfig json.RawMessage) (*WithRDB, error) {
	cfg, err := database.ConfigFromJson(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	return NewWithRDB(cfg)
}

func NewWithRDB(databaseConfig *database.Config) (*WithRDB, error) {
	db, err := gorm.Open(postgres.Open(databaseConfig.ToConnectStr()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &WithRDB{
		BaseRDB: db,
	}, nil
}

// WithDB defines a struct that has capability to access the database.
type WithRDB struct {
	BaseRDB *gorm.DB
}

func (s *WithRDB) ConfigureReaderDatabase(*WithRDB) {
	panic("ConfigureReaderDatabase must be implemented")
}

// DB returns an usable DB connection
func (s *WithRDB) RDB(ctx context.Context) *gorm.DB {
	if s.BaseRDB == nil {
		return nil
	}

	return s.BaseRDB.Session(&gorm.Session{})
}
