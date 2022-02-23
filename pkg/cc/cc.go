package cc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/cc/config"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	"github.com/aserto-dev/aserto-idp/pkg/provider/finder"
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// CC contains dependencies that are cross cutting and are needed in most
// of the providers that make up this application
type CC struct {
	Context     context.Context
	Config      *config.Config
	Log         *zerolog.Logger
	UI          *clui.UI
	Retriever   retriever.Retriever
	pluginsInfo *retriever.PluginsInfo
	providers   map[string]provider.Provider
}

func (c *CC) SetLogger(w io.Writer) {
	log.SetOutput(w)
}

func (c *CC) GetLatestVersion(pluginName string) (string, error) {
	return c.pluginsInfo.LatestVersion(pluginName)
}

func (c *CC) GetRemotePluginsInfo() (map[string]map[string][]string, error) {
	return c.pluginsInfo.GetInfo()
}

// ProviderExists returns true if the provider has already been added to the context
func (c *CC) ProviderExists(name string) bool {
	_, ok := c.providers[name]
	return ok
}

// AddProvider to the context
func (c *CC) AddProvider(prov provider.Provider) error {
	if c.providers == nil {
		c.providers = make(map[string]provider.Provider)
	}
	provName := prov.GetName()
	if c.ProviderExists(provName) {
		return fmt.Errorf("provider %s has already been added", provName)
	}

	c.providers[prov.GetName()] = prov
	return nil
}

func (c *CC) GetProviders() map[string]provider.Provider {
	return c.providers
}

// GetProvider with the given name
func (c *CC) GetProvider(name string) provider.Provider {
	return c.providers[name]
}

// Dispose all the resources. This can be called any number of times
func (c *CC) Dispose() {
	for _, provider := range c.providers {
		provider.Kill()
	}
}

func (c *CC) ConnectRetriever() error {
	err := c.Retriever.Connect()
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig loads the plugin and logger config from a configuration file
func (c *CC) LoadConfig(path string) error {
	cfg, err := config.NewConfig(path, c.Log)
	if err != nil {
		return errors.Wrap(err, "error while loading configuration")
	}
	c.Log.Debug().Msgf("using config file %s", path)
	c.Config = cfg

	if cfg.Logging != nil && c.Log.GetLevel() == zerolog.ErrorLevel {
		newLogger, err := logger.NewLogger(os.Stdout, os.Stderr, cfg.Logging)
		if err != nil {
			c.Log.Warn().Msgf("failed to load logger from config file '%s'", err.Error())
		} else {
			c.Log = newLogger
		}
	}
	return nil
}

func (c *CC) LoadProviders() error {
	envFinder := finder.NewHomeDir()

	pluginPaths, err := envFinder.Find()
	if err != nil {
		return err
	}
	for _, pluginPath := range pluginPaths {
		idpProvider := provider.NewIDPProvider(c.Log, pluginPath)

		if c.ProviderExists(idpProvider.GetName()) {
			c.Log.Printf("Plugin %s has already been loaded from %s. Ignoring %s", idpProvider.GetName(), idpProvider.GetPath(), pluginPath)
			continue
		}

		err = c.AddProvider(idpProvider)
		if err != nil {
			c.Log.Printf("could not add provider %s, error: %s", idpProvider.GetName(), err.Error())
			continue
		}

	}

	return nil
}

func GetLogLevel() zerolog.Level {
	logLevel := zerolog.FatalLevel

	for _, arg := range os.Args {
		if strings.HasPrefix(strings.ToLower(arg), "--verbosity=") {
			intLevel, err := strconv.Atoi(strings.Split(arg, "=")[1])
			if err != nil {
				break
			}
			switch intLevel {
			case 1:
				logLevel = zerolog.ErrorLevel
			case 2:
				logLevel = zerolog.InfoLevel
			case 3:
				logLevel = zerolog.DebugLevel
			case 4:
				logLevel = zerolog.TraceLevel
			}

		}
		switch arg {
		case "-v":
			logLevel = zerolog.ErrorLevel
		case "-vv":
			logLevel = zerolog.InfoLevel
		case "-vvv":
			logLevel = zerolog.DebugLevel
		case "-vvvv":
			logLevel = zerolog.TraceLevel
		}

	}

	return logLevel
}
