package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/provider/finder"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/pkg/errors"
)

func main() {
	c := cc.New()

	cli := cmd.CLI{}

	envFinder := finder.NewEnvironment()

	pluginOptions, err := cmd.LoadPlugins(envFinder)
	if err != nil {
		log.Fatal(err.Error())
	}

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
		kong.Vars{"defaultEnv": x.DefaultEnv},
	}
	options = append(options, pluginOptions...)

	ctx := kong.Parse(&cli, options...)

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

	err = ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
