package dbtest

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/database"

	// used because the source of the migration is a file.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultDBName   = "postgres"
	defaultSuffix   = "_test"
)

// Config contains a set of properties that can be redefined to accommodate
// different test environments.
var Config = struct {
	// Username is the database username used to open the connection.
	Username string

	// Password is the database password used to open the connection.
	Password string

	// RootDBName is the database name to connect to. This is not the same dbname
	// that the tests will run within.
	// This is only used to create the initial connection that will then create
	// the database for the tests (based on the argument given to `DB()`).
	RootDBName string

	// DBSuffix is a string that is appended to the database name to
	// distinguish it from the "non-testing" version.
	// e.g. if DBSuffix is `_test`, and the name given to DB() is `myservice`,
	// for testing purposes, this package will create and use the `myservice_test`
	// database instead.
	// Empty DBSuffix will cause this package to use the raw name provided in DB()
	// for testing too.
	DBSuffix string
}{
	Username:   defaultUser,
	Password:   defaultPassword,
	RootDBName: defaultDBName,
	DBSuffix:   defaultSuffix,
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
func DB(t *testing.T, migrationsPath, name string) (*sql.DB, func()) {

	name = name + Config.DBSuffix

	var db *sql.DB

	if !isDBInitialized(name) {
		db = initializeDB(t, migrationsPath, name)

	} else {
		var err error
		connStr := connectionString(defaultUser, defaultPassword, name)
		db, err = sql.Open("postgres", connStr)
		require.NoError(t, err, "failed to open DB connection")

	}

	resetDB(t, db, name)

	return db, func() { defer Close(t, db) }
}

// Close can be used to close a database connection.
func Close(t *testing.T, db *sql.DB) {
	require.NoError(t, db.Close())
}

func isDBInitialized(name string) bool {
	_, ok := initializedDBs[name]
	return ok
}

func initializeDB(t *testing.T, migrationsPath, name string) *sql.DB {
	connStr := connectionString(Config.Username, Config.Password, Config.RootDBName)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to open DB connection")

	// Intentionally ignore errors here
	// Ideally we want to drop and recreate on the first run, but we don't
	// want to fail in case someone is connected to the db for example.
	_, _ = db.Exec("DROP DATABASE IF EXISTS " + name)
	_, _ = db.Exec("CREATE DATABASE " + name)

	initializedDBs[name] = struct{}{}

	require.NoError(t, db.Close())

	connStr = connectionString(defaultUser, defaultPassword, name)
	db, err = sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to open DB connection")

	if migrationsPath != "" {
		migrateDB(t, db, migrationsPath)
	}

	return db
}

func resetDB(t *testing.T, db *sql.DB, name string) {
	rows, err := db.Query("SELECT relname FROM pg_stat_user_tables")
	require.NoError(t, err)

	defer rows.Close()

	for rows.Next() {
		var relname string

		require.NoError(t, rows.Scan(&relname))

		// Do not clean up the schema migrations.
		if relname == "schema_migrations" {
			continue
		}

		_, err := db.Exec("TRUNCATE " + relname + " CASCADE")
		require.NoError(t, err)
	}
}

func connectionString(user, password, dbname string) string {
	return fmt.Sprintf("host=127.0.0.1 port=5432 sslmode=disable user=%s password=%s dbname=%s", user, password, dbname)
}

func migrateDB(t *testing.T, db *sql.DB, migrationsPath string) {
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
