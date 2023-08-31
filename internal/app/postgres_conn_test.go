package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectToPostgresInitDB(t *testing.T) {
	dsn := "postgres://postgres:postgres@localhost:11543/postgres?sslmode=disable"

	db, err := ConnectToPostgres(dsn)
	require.NoError(t, err)
	require.NotNil(t, db)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
}
