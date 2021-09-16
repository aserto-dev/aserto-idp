package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var log zerolog.Logger

func getConfigFile(content string, t *testing.T) string {
	tmpDir := t.TempDir()

	configFilePath := filepath.Join(tmpDir, fmt.Sprintf("%s.yaml", t.Name()))
	configFile, err := os.Create(configFilePath)
	if err != nil {
		t.Error("cannot create configuration file")
	}
	defer configFile.Close()

	_, err = configFile.WriteString(content)
	if err != nil {
		t.Error(err)
	}
	return configFilePath
}

func TestConfigDoesNotExist(t *testing.T) {
	assert := require.New(t)
	configPath := filepath.Join(t.TempDir(), "not_exist.yaml")

	cfg, err := NewConfig(configPath, &log)
	assert.NotNil(err)
	assert.Nil(cfg)

}

func TestConfigFromFile(t *testing.T) {
	assert := require.New(t)
	content := `
---
logging:
  log_level: trace

plugins:
  json:
    file: user_import.json
  aserto:
    tenant: "00000000-0000-0000-0000-000000000000"
    authorizer: "authorizer.prod.aserto.com:8443"
    api_key: "0000000000000000000000000000000000000000000000000000"
`

	configFilePath := getConfigFile(content, t)
	cfg, err := NewConfig(configFilePath, &log)

	assert.Nil(err)
	assert.NotNil(cfg.Plugins["json"])
	assert.NotNil(cfg.Plugins["aserto"])
	assert.Equal(zerolog.LevelTraceValue, cfg.Logging.LogLevel)
	assert.Nil(cfg.Plugins["ion"])

	jsonSettings := cfg.Plugins["json"]
	assert.Equal("user_import.json", jsonSettings["file"])
	assert.Nil(jsonSettings["nosetting"])

	asertoSettings := cfg.Plugins["aserto"]
	assert.Equal("00000000-0000-0000-0000-000000000000", asertoSettings["tenant"])
	assert.Equal("authorizer.prod.aserto.com:8443", asertoSettings["authorizer"])
	assert.Equal("0000000000000000000000000000000000000000000000000000", asertoSettings["api_key"])

}

func TestConfigFromFileNoPlugins(t *testing.T) {
	assert := require.New(t)
	content := `
---
logging:
  log_level: error
`

	configFilePath := getConfigFile(content, t)
	cfg, err := NewConfig(configFilePath, &log)

	assert.Nil(err)
	assert.Nil(cfg.Plugins["aserto"])
	assert.Equal(zerolog.ErrorLevel, cfg.Logging.LogLevelParsed)
}

func TestConfigFromFileNoPLogging(t *testing.T) {
	assert := require.New(t)
	content := `
---
plugins:
  json:
    file: user_import.json
`

	configFilePath := getConfigFile(content, t)
	cfg, err := NewConfig(configFilePath, &log)

	assert.Nil(err)
	assert.Nil(cfg.Logging)
}
