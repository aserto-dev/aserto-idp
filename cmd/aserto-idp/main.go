package main

import (
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	"github.com/aserto-dev/aserto-idp/pkg/provider/finder"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/logger"
)

func main() {
	cfg := &logger.Config{}
	cfg.LogLevelParsed = cc.GetLogLevel()

	c, err := cc.BuildCC(os.Stdout, os.Stderr, os.Stdout, cfg)
	if err != nil {
		log.Fatalf("failed to build application: %s", err.Error())
	}

	err = appStart(c)

	if err != nil {
		c.UI.Problem().Msg(err.Error())
		os.Exit(1)
	}
}

func appStart(c *cc.CC) error {
	cli := cmd.CLI{}

	envFinder := finder.NewHomeDir()

	pluginPaths, err := envFinder.Find()
	if err != nil {
		return err
	}

	defer func() {
		c.Dispose()
		c.Retriever.Disconnect()
	}()

	options := []kong.Option{
		kong.Name(x.AppName),
		kong.Exit(func(exitCode int) {
			c.Dispose()
			os.Exit(exitCode)
		}),
		kong.Description(x.AppDescription),
		kong.UsageOnError(),

		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             false,
			Compact:             true,
			Tree:                false,
			FlagsLast:           true,
			Indenter:            kong.SpaceIndenter,
			NoExpandSubcommands: false,
		}),
		kong.Bind(c),
		kong.Vars{"defaultEnv": x.DefaultEnv},
	}

	for _, pluginPath := range pluginPaths {
		idpProvider := provider.NewIDPProvider(c.Log, pluginPath)

		if c.ProviderExists(idpProvider.GetName()) {
			log.Printf("Plugin %s has already been loaded from %s. Ignoring %s", idpProvider.GetName(), idpProvider.GetPath(), pluginPath)
			continue
		}

		plugin, err := cmd.NewPlugin(idpProvider, c)
		if err != nil {
			log.Printf("could not create new plugin %s, error: %s", idpProvider.GetName(), err.Error())
			continue
		}

		err = c.AddProvider(idpProvider)
		if err != nil {
			log.Printf("could not add provider %s, error: %s", idpProvider.GetName(), err.Error())
			continue
		}

		cli.Plugins = append(cli.Plugins, plugin.Plugins...)

	}
	ctx := kong.Parse(&cli, options...)

	err = c.LoadConfig(strings.TrimSpace(cli.Config))
	if err != nil {
		return err
	}

	// TODO add config option for custom package repo
	err = c.ConnectRetriever()
	if err != nil {
		return err
	}

	err = ctx.Run(c)

	if err != nil {
		return err
	}

	return nil
}
