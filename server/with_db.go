package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tscolari/servicetools/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewWithDBFromJSONConfig(jsonConfig json.RawMessage) (*WithDB, error) {
	cfg, err := database.ConfigFromJson(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	return NewWithDB(cfg)
}

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
type WithDB struct {
	BaseDB *gorm.DB
}

func (s *WithDB) ConfigureDatabase(*WithDB) {
	panic("ConfigureDatabase must be implemented")
}

// DB returns an usable DB connection
func (s *WithDB) DB(ctx context.Context) *gorm.DB {
	if s.BaseDB == nil {
		return nil
	}

	return s.BaseDB.Session(&gorm.Session{})
}
