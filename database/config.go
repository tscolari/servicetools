package database

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	envHostname = "HOSTNAME"
	envPort     = "PORT"
	envUsername = "USERNAME"
	envPassword = "PASSWORD"
	envDBName   = "NAME"
	envSSLMode  = "SSLMODE"
)

var ErrNoEnvConfiguration = errors.New("no database configuration available on environmental variables")

// Config defines the basic configuration for connecting to the database.
type Config struct {
	Hostname string `json:"hostname,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"db_name,omitempty"`
	SSLMode  bool   `json:"ssl_mode,omitempty"`
}

// ConfigFromEnv loads the database configuration from env variables
// using the given prefix and the `_` separator.
func ConfigFromEnv(prefix string) (*Config, error) {
	config := Config{
		Hostname: os.Getenv(fmt.Sprintf("%s_%s", prefix, envHostname)),
		Username: os.Getenv(fmt.Sprintf("%s_%s", prefix, envUsername)),
		Password: os.Getenv(fmt.Sprintf("%s_%s", prefix, envPassword)),
		DBName:   os.Getenv(fmt.Sprintf("%s_%s", prefix, envDBName)),
		SSLMode:  os.Getenv(fmt.Sprintf("%s_%s", prefix, envSSLMode)) == "true",
	}

	// Check if we actually have any value from the env up to this point.
	if reflect.DeepEqual(config, Config{}) {
		return nil, ErrNoEnvConfiguration
	}

	port, err := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_%s", prefix, envPort)))
	if err != nil {
		return nil, err
	}

	config.Port = port

	return &config, nil
}

// ToConnectStr returns a connection string using the config's values.
func (c Config) ToConnectStr() string {
	return fmt.Sprintf(
		"host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
		c.Hostname,
		c.Port,
		c.sslMode(),
		c.Username,
		c.Password,
		c.DBName,
	)
}

func (c Config) sslMode() string {
	if c.SSLMode {
		return "enabled"
	}

	return "disable"
}
