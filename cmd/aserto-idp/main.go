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

	pluginsMap := make(map[string]*cmd.Plugin)

	options := []kong.Option{
		kong.Name(x.AppName),
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
		idpProvider := provider.NewIDPProvider(pluginPath)

		if path, ok := pluginsMap[idpProvider.GetName()]; ok {
			log.Printf("Plugin %s has already been loaded from %s. Ignoring %s", idpProvider.GetName(), path, pluginPath)
			continue
		}
		plugin, err := cmd.NewPlugin(idpProvider)
		if err != nil {
			log.Fatal(err.Error())
		}

		pluginsMap[idpProvider.GetName()] = plugin

		if plugin.Name == x.DefaultPluginName {
			client, err := idpProvider.PluginClient()
			if err != nil {
				log.Fatal(err.Error())
			}
			c.DefaultIDPClient = client
			cli.Plugins = append(cli.Plugins, plugin.Plugins...)
		} else {
			dynamicCommand := kong.DynamicCommand(plugin.Name, plugin.Description, "Plugins", plugin)
			options = append(options, dynamicCommand)
			c.IDPClients[plugin.Name], err = idpProvider.PluginClient()
			if err != nil {
				log.Fatal(err.Error())
			}
		}

	}

	ctx := kong.Parse(&cli, options...)

	if cli.Debug {
		c.SetLogger(os.Stderr)
	}

	err = ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
