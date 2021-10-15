package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/version"
	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type CLI struct {
	Config    string     `short:"c" type:"path" help:"Path to the config file. Any argument provided to the CLI will take precedence."`
	Delete    DeleteCmd  `cmd:"" help:"delete user ids from an user-provider idp"`
	Exec      ExecCmd    `cmd:"" help:"import users from an user-provided idp to another user-provided idp"`
	Version   VersionCmd `cmd:"" help:"version information"`
	Verbosity int        `short:"v" type:"counter" help:"Use to increase output verbosity."`
	kong.Plugins
}

func (cmd *CLI) Run(c *cc.CC) error {
	return nil
}

type VersionCmd struct {
}

func (cmd *VersionCmd) Run(c *cc.CC) error {
	fmt.Printf("%s - %s (%s)\n",
		x.AppName,
		version.GetVersionString(),
		x.AppVersionTag,
	)
	return nil
}
