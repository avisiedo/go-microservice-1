package datastore

import (
	"log/slog"
	"testing"

	"github.com/avisiedo/go-microservice-1/internal/config"
	"github.com/avisiedo/go-microservice-1/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperConfig() *config.Config {
	return &config.Config{
		Database: config.Database{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			Name:     "dbname",
		},
	}
}

func TestGetURL(t *testing.T) {
	result := getURL(&config.Config{})
	assert.Equal(t, "user= password= host= port=0 dbname= sslmode=disable", result)

	result = getURL(helperConfig())
	assert.Equal(t, "user=user password=password host=localhost port=5432 dbname=dbname sslmode=disable", result)

	result = getURL(&config.Config{
		Database: config.Database{
			Host:       "localhost",
			Port:       5432,
			User:       "user",
			Password:   "password",
			Name:       "dbname",
			CACertPath: "/tmp/ca.cert",
		},
	})
	assert.Equal(t, "user=user password=password host=localhost port=5432 dbname=dbname sslmode=verify-full sslrootcert=/tmp/ca.cert", result)
}

// This test require the infrastructure to be started
func TestNewDB(t *testing.T) {
	t.Skip("require a database up and running for unit testing")
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		_ = NewDB(nil, slog.Default())
	})

	cfg := config.Get()
	require.NotNil(t, cfg)
	oldPort := cfg.Database.Port
	cfg.Database.Port = 2345
	db := NewDB(cfg, nil)
	cfg.Database.Port = oldPort
	require.Nil(t, db)

	assert.NotPanics(t, func() {
		_ = NewDB(cfg, nil)
	})

	db = NewDB(cfg, slog.Default())
	require.NotNil(t, db)
}

func TestClose(t *testing.T) {
	assert.NotPanics(t, func() {
		Close(nil)
	})

	mock, db, err := test.NewSqlMock(nil)
	require.NoError(t, err)
	require.NotNil(t, db)
	mock.ExpectClose()
	assert.NotPanics(t, func() {
		Close(db)
	})
}

func TestDBMigrationPath(t *testing.T) {
	oldDBMigrationPath := dbMigrationPath
	dbMigrationPath = ""
	assert.Equal(t, defaultDBMigrationPath, DBMigrationPath())
	dbMigrationPath = oldDBMigrationPath

	assert.Equal(t, dbMigrationPath, DBMigrationPath())
}

func TestSetDBMigrationPath(t *testing.T) {
	oldDBMigrationPath := dbMigrationPath
	SetDBMigrationPath("/tmp")
	assert.Equal(t, "/tmp", DBMigrationPath())
	SetDBMigrationPath(oldDBMigrationPath)
}
