package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"

	"github.com/tscolari/servicetools/database"
)

// CanMigrate injects the "migrate" subcommand to another command.
func CanMigrate(rootCmd *cobra.Command) {
	rootCmd.AddCommand(migrateCmd)
}

func init() {
	// path should point to a folder migration files.
	migrateCmd.PersistentFlags().StringVarP(&migratePath, "path", "p", "./migrations", "path to all migrations")
	migrateCmd.PersistentFlags().StringVarP(&migrateEnvPrefix, "db-env-prefix", "e", "DATABASE", "prefix for all DB env variables")
}

var (
	migratePath      string
	migrateEnvPrefix string
)

// migrateCmd performs the database migration for the given path.
// The connection to the database will use the given `db-env-prefix` for:
// "_HOSTNAME", "_PORT", "_USERNAME", "_PASSWORD", "_NAME" and "_SSLMODE".
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database with the given migrations",
	RunE: func(cmd *cobra.Command, args []string) error {

		migrationStat, err := os.Stat(migratePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to identify the migrations path: %v", err)
			return err
		}

		if !migrationStat.IsDir() {
			fmt.Fprintf(os.Stderr, "the given migrations path is not a directory")
			return errors.New("the given migration path is not a directory")
		}

		var dbConfig *database.Config

		dbConfig, err = database.ConfigFromEnv(migrateEnvPrefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load configuration from env: %v\n", err)
			return err
		}

		db, err := sql.Open("postgres", dbConfig.ToConnectStr())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
			return err
		}

		if err := db.Ping(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to check connection to database: %v\n", err)
			return err
		}

		if err := database.Migrate(db, migratePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate the database: %v\n", err)
			return err
		}

		return nil
	},
}
