package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/pkg/errors"
)

func main() {
	c := cc.New()
	plugins := cmd.FindPlugins()

	cli := cmd.CLI{}
	for _, pl := range plugins {
		h, err := cmd.GetPluginHelp(pl)
		if err == nil {
			cli.Export.Plugins = append(cli.Export.Plugins, h)
		}
	}

	ctx := kong.Parse(&cli,
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
		kong.Vars{"defaultEnv": x.DefaultEnv},
		kong.Bind(&cli),
	)

	if cli.Debug {
		c.SetLogger(os.Stderr)
	}

	if cli.APIKey != "" {
		c.SetAPIKey(cli.APIKey)
	}

	if err := c.SetEnv(cli.EnvOverride); err != nil {
		ctx.FatalIfErrorf(errors.Wrapf(err, "set environment [%s]", cli.EnvOverride))
	}

	if cli.TenantOverride != "" {
		c.Override(x.TenantIDOverride, cli.TenantOverride)
	}

	if cli.AuthorizerOverride != "" {
		c.Override(x.AuthorizerOverride, cli.AuthorizerOverride)
	}

	if cli.Provider != "" {
		plugin, _ := cmd.LoadPlugin(cli.Provider, plugins[cli.Provider])
		c.SetPlugin(plugin)
	}

	err := ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
