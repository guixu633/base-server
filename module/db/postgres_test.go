package db

import (
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func getDB(t *testing.T) *Database {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	db, err := NewDB(cfg.Postgres)
	assert.NoError(t, err)
	return db
}
