package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/pkg/errors"
)

type VetaCMD struct {
}

func (v *VetaCMD) Run() error {
	return nil
}

// type plug struct {
// 	Vata  VetaCMD `cmd`
// 	param string
// }

//plugJson := `{"cmd":"mycmd"}`

func main() {
	c := cc.New()
	plugins := cmd.FindPlugins()

	cli := cmd.CLI{}
	for _, plugin := range plugins {
		//help, err := cmd.GetPluginHelp(plugin)
		fmt.Println("!!!!!!!! -- " + plugin)
		// varPlugMap := make(map[string]VetaCMD)
		// varPlugMap["value"] = VetaCMD{}
		//if err == nil {
		//cli.Plugins = append(cli.Plugins, &plug{})
		//}
	}

	// var beforeApplt kong.BeforeApply = nil
	// var object kong.BeforeApply = func() {
	// 	fmt.Println("eureka!")
	// }
	opt := kong.DynamicCommand("ion", "this is a json plugin", "Plugins", &VetaCMD{})

	ctx := kong.Parse(&cli,
		kong.Name(x.AppName),
		kong.Description(x.AppDescription),
		kong.UsageOnError(),
		//kong.NamedMapper("moo", testMooMapper{}),
		opt,
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
		//kong.Bind(&cli),
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
		plugin, _ := cmd.LoadPlugin(string(cli.Provider), plugins[string(cli.Provider)])
		c.SetPlugin(plugin)
	}

	err := ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
