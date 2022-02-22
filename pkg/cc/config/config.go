package config

import (
	"os"
	"strings"

	"github.com/aserto-dev/go-utils/logger"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Logging *logger.Config                    `json:"logging"`
	Plugins map[string]map[string]interface{} `json:"plugins"`
}

func NewEmptyConfig() *Config {
	return &Config{}
}

// Loads the config from a file.
func NewConfig(configPath string, log *zerolog.Logger) (*Config, error) {
	if configPath == "" {
		return &Config{}, nil
	}
	configLogger := log.With().Str("component", "config").Logger()
	log = &configLogger

	exists, err := fileExists(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine if the config file %s exists", configPath)
	}

	if !exists {
		return nil, errors.Errorf("Config file %s does not exist", configPath)
	}

	v := viper.New()
	v.SetConfigFile("yaml")
	v.AddConfigPath(".")
	v.SetConfigFile(configPath)
	v.SetEnvPrefix("ASERTO_IDP_CLI")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err = v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open config file '%s'", configPath)
	}
	v.AutomaticEnv()

	cfg := new(Config)
	err = v.UnmarshalExact(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config file '%s'", configPath)
	}

	if cfg.Logging != nil {
		cfg.Logging.LogLevelParsed, err = zerolog.ParseLevel(cfg.Logging.LogLevel)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse log level %s", cfg.Logging.LogLevel)
		}
	}

	return trimPluginName(cfg), nil
}

func trimPluginName(cfg *Config) *Config {
	newCfg := &Config{Logging: cfg.Logging}
	newCfg.Plugins = make(map[string]map[string]interface{})
	for plugin, options := range cfg.Plugins {
		newCfg.Plugins[plugin] = make(map[string]interface{})
		for option, value := range options {
			trimmedOptionName := strings.TrimPrefix(option, plugin+"-")
			newCfg.Plugins[plugin][trimmedOptionName] = value
		}
	}
	return newCfg
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "failed to stat file '%s'", path)
	}
}
