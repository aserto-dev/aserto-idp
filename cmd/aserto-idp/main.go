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
)

func main() {
	c := cc.New()

	cli := cmd.CLI{}

	envFinder := finder.NewEnvironment()

	pluginPaths, err := envFinder.Find()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		c.Dispose()
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
			log.Fatal(err.Error())
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
		c.Log.Fatal().Msg(err.Error())
	}

	err = ctx.Run(c)

	if err != nil {
		c.Ui.Problem().Msg(err.Error())
	}
}
