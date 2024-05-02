package dbtest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// used because the source of the migration is a file.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultDBName   = "postgres"
	defaultSuffix   = "_test"
)

// Config contains a set of properties that can be redefined to accomodate
// different test environments.
var Config = struct {
	// Username is the database username used to open the connection.
	Username string

	// Password is the database password used to open the connection.
	Password string

	// DBName is the database name to connect to. This is not the same dbname
	// that the tests will run within.
	// This is only used to create the initial connection that will then create
	// the database for the tests (based on the argument given to `DB()`).
	DBName string

	// DBSuffix is a string that is appended to the database name to
	// distinguish it from the "non-testing" version.
	DBSuffix string
}{
	Username: defaultUser,
	Password: defaultPassword,
	DBName:   defaultDBName,
	DBSuffix: defaultSuffix,
}

var initializedDBs map[string]struct{}

func init() {
	initializedDBs = map[string]struct{}{}
}

// DB is meant to be used in tests.
// It will take a migrations path and a database name to be used.
// The first time it gets called, it will ensure that the database
// exists (it will be dropped and recreated if possible) and migrate
// the database.
// On every call it will truncate all the tables (except the schema one)
// to ensure that there is no data contamination.
// It will always append `Config.DBSuffix` to the given database name,
// to have no name appended, set `dbtest.Config.DBSuffix = ""`.
//
// Because DB will reset the database on every call, it's not safe for
// this to be used in parallel tests, unless they are using different
// database names.
func DB(t *testing.T, migrationsPath, name string) (*gorm.DB, func()) {

	name = name + Config.DBSuffix

	var db *gorm.DB

	if !isDBInitialized(name) {
		db = initializeDB(t, migrationsPath, name)

	} else {
		var err error
		connStr := connectionString(defaultUser, defaultPassword, name)
		db, err = gorm.Open(postgres.Open(connStr))
		require.NoError(t, err, "failed to open DB connection")

	}

	resetDB(t, db, name)

	return db, func() { defer Close(t, db) }
}

// Close can be used to close a database connection.
func Close(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
}

func isDBInitialized(name string) bool {
	_, ok := initializedDBs[name]
	return ok
}

func initializeDB(t *testing.T, migrationsPath, name string) *gorm.DB {
	connStr := connectionString(Config.Username, Config.Password, Config.DBName)
	db, err := gorm.Open(postgres.Open(connStr))
	require.NoError(t, err, "failed to open DB connection")

	defer Close(t, db)

	// Intentionally ignore errors here
	// Ideally we want to drop and recreate on the first run, but we don't
	// want to fail in case someone is connected to the db for example.
	_ = db.Exec("DROP DATABASE IF EXISTS " + name).Error
	_ = db.Exec("CREATE DATABASE " + name).Error

	initializedDBs[name] = struct{}{}

	connStr = connectionString(defaultUser, defaultPassword, name)
	db, err = gorm.Open(postgres.Open(connStr))
	require.NoError(t, err, "failed to open DB connection")

	if migrationsPath != "" {
		migrateDB(t, db, migrationsPath)
	}

	return db
}

func resetDB(t *testing.T, db *gorm.DB, name string) {
	rows, err := db.Table("pg_stat_user_tables").Rows()
	require.NoError(t, err)

	for rows.Next() {
		table := struct {
			Relname    string
			Schemaname string
		}{}

		require.NoError(t, db.ScanRows(rows, &table))

		// Do not clean up the schema migrations.
		if table.Relname == "schema_migrations" {
			continue
		}

		require.NoError(t, db.Exec("TRUNCATE "+table.Relname+" CASCADE").Error)
	}
}

func connectionString(user, password, dbname string) string {
	return fmt.Sprintf("host=127.0.0.1 port=5432 sslmode=disable user=%s password=%s dbname=%s", user, password, dbname)
}

func migrateDB(t *testing.T, db *gorm.DB, migrationsPath string) {
	dir, err := os.Getwd()
	require.NoError(t, err, "failed to get current directory")
	baseDir := dir
	defer func() {
		_ = os.Chdir(dir)
	}()

	var path string
	for {
		path = filepath.Join(baseDir, migrationsPath)
		if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
			require.NoError(t, os.Chdir(".."), "failed to change paths")
			baseDir, err = os.Getwd()
			require.NoError(t, err, "failed to get current directory")

			if baseDir == "/" {
				require.Fail(t, "failed to find migrations path")
			}
		} else {
			break
		}
	}

	require.NoError(t, database.Migrate(db, path))
}
