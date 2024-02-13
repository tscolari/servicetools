package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/database"
)

type testWithDB struct {
	*WithDB
}

func (t *testWithDB) ConfigureDatabase(w *WithDB) {
	t.WithDB = w
}

func Test_WithDB(t *testing.T) {
	config := database.Config{
		Hostname: "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  false,
	}

	withDB, err := NewWithDB(&config)
	require.NoError(t, err)

	testObj := &testWithDB{}
	testObj.ConfigureDatabase(withDB)

	db := testObj.DB(context.Background())
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
