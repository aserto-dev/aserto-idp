package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/version"
	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type CLI struct {
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
	fmt.Fprintf(c.Log, "%s - %s (%s)\n",
		x.AppName,
		version.GetInfo().String(),
		x.AppVersionTag,
	)
	return nil
}
