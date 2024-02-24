package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/database"
)

type testWithRDB struct {
	*WithRDB
}

func (t *testWithRDB) ConfigureReaderDatabase(w *WithRDB) {
	t.WithRDB = w
}

func Test_WithRDB(t *testing.T) {
	config := database.Config{
		Hostname: "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  false,
	}

	withRDB, err := NewWithRDB(&config)
	require.NoError(t, err)

	testObj := &testWithRDB{}
	testObj.ConfigureReaderDatabase(withRDB)

	db := testObj.RDB(context.Background())
	require.NotNil(t, db)

	testQuery := `SELECT * FROM information_schema.tables`
	require.NoError(t, db.Exec(testQuery).Error)

	t.Run("when the connection is invalid", func(t *testing.T) {
		config := database.Config{
			Hostname: "localhost",
			Port:     5432,
			Username: "wrongUser",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  false,
		}

		_, err := NewWithDB(&config)
		require.Error(t, err)
	})
}

