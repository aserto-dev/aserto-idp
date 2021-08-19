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

	// client := plugin.NewClient(&plugin.ClientConfig{
	// 	HandshakeConfig: shared.Handshake,
	// 	Plugins:         shared.PluginMap,
	// 	Cmd:             exec.Command("/home/florin_aserto_com/_code/aserto-idp/main"),
	// 	AllowedProtocols: []plugin.Protocol{
	// 		plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	// })
	// defer client.Kill()

	// rpcClient, err := client.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // Request the plugin
	// raw, err := rpcClient.Dispense("provider")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// p := raw.(shared.Provider)
	// pluginHelp, err := p.Help()
	// var res shared.HelpMessage

	// structpbconv.Convert(pluginHelp.HelpStruct, &res)

	plugins := cmd.FindPlugins()
	// help := []*shared.HelpMessage{}

	// var help *shared.HelpMessage
	// var err error

	cli := cmd.CLI{}
	for _, pl := range plugins {
		h, err := cmd.GetPluginHelp(pl)
		if err == nil {
			cli.Import.Plugins = append(cli.Import.Plugins, h)
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
