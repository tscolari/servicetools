package cmd

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
	"github.com/tscolari/servicetools/database"
)

func CanMigrate(rootCmd *cobra.Command) {
	rootCmd.AddCommand(migrateCmd)
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&migrationsPath, "path", "p", "./migrations", "path to migrations")
}

var (
	migrationsPath string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database with the given migrations",
	RunE: func(cmd *cobra.Command, args []string) error {

		migrationStat, err := os.Stat(migrationsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to identify the migrations path: %v", err)
			return err
		}

		if !migrationStat.IsDir() {
			fmt.Fprintf(os.Stderr, "the given migrations path is not a directory")
			return errors.New("the given migration path is not a directory")
		}

		var dbConfig *database.Config

		dbConfig, err = database.ConfigFromEnv("DATABASE")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load configuration from env: %v\n", err)
			return err
		}

		db, err := gorm.Open(postgres.Open(dbConfig.ToConnectStr()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to database: %w\n", err)
			return err
		}

		if err := database.Migrate(db, migrationsPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate the database: %w\n", err)
			return err
		}

		return nil
	},
}
