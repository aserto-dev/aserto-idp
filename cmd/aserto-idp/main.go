package main

import (
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/provider/finder"
	"github.com/aserto-dev/aserto-idp/pkg/x"
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
		kong.Bind(c),
		kong.Vars{"defaultEnv": x.DefaultEnv},
	}
	options = append(options, pluginOptions...)

	ctx := kong.Parse(&cli, options...)
	c.Command = strings.Fields(ctx.Command())[0]

	if cli.Debug {
		c.SetLogger(os.Stderr)
	}

	err = cmd.SetPluginContext(c, envFinder)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
