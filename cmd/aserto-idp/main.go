package main

import (
	"log"
	"os"

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

		if plugin.Name == x.DefaultPluginName {
			err := c.SetDefaultProvider(idpProvider)
			if err != nil {
				log.Fatal(err.Error())
			}
			cli.Plugins = append(cli.Plugins, plugin.Plugins...)
		} else {
			err = c.AddProvider(idpProvider)
			if err != nil {
				log.Printf("could not add provider %s, error: %s", idpProvider.GetName(), err.Error())
				continue
			}
			dynamicCommand := kong.DynamicCommand(plugin.Name, plugin.Description, "Plugins", plugin)
			options = append(options, dynamicCommand)

		}

	}

	ctx := kong.Parse(&cli, options...)

	err = ctx.Run(c)

	ctx.FatalIfErrorf(err)
}
