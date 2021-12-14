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
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/aserto-dev/aserto-idp/pkg/x"
)

func main() {
	c := cc.New()

	err := appStart(c)

	if err != nil {
		log.Fatal(err.Error())
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
		kong.Description(constructDescription(c.Retriever)),
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

		err = c.AddProvider(idpProvider)
		if err != nil {
			log.Printf("could not add provider %s, error: %s", idpProvider.GetName(), err.Error())
			continue
		}

		plugin, err := cmd.NewPlugin(idpProvider, c)
		if err != nil {
			return err
		}

		cli.Plugins = append(cli.Plugins, plugin.Plugins...)

	}
	ctx := kong.Parse(&cli, options...)

	err = c.LoadConfig(strings.TrimSpace(cli.Config))
	if err != nil {
		return err
	}

	//TODO add config option for custom package repo
	err = c.ConnectRetriever()
	if err != nil {
		return err
	}

	err = ctx.Run(c)

	if err != nil {
		c.Ui.Problem().Msg(err.Error())
		return err
	}

	return nil
}

func constructDescription(ghcr retriever.Retriever) string {

	plugins := retriever.PluginVersions(ghcr)
	if len(plugins) == 0 {
		return x.AppDescription
	}

	header := x.AppDescription + "\n " + "Plugins available to download:"

	for key := range plugins {
		header = header + "\n " + key
	}

	return header
}
