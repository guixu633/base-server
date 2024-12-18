package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../../config.toml")
	assert.NoError(t, err)
	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(jsonCfg))
}
