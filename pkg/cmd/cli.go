package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/version"
	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type CLI struct {
	Version            VersionCmd `cmd:"" help:"version information"`
	Import             ImportCmd  `cmd:"" help:"import users"`
	Provider           string     `required:"" help:"load users provider (json | auth0)" enum:"json,auth0"`
	Verbose            bool       `name:"verbose" help:"verbose output"`
	Debug              bool       `name:"debug" env:"ASERTO_DEBUG" help:"enable debug logging"`
	AuthorizerOverride string     `name:"authorizer" env:"ASERTO_AUTHORIZER" help:"authorizer override"`
	TenantOverride     string     `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
	EnvOverride        string     `name:"env" default:"${defaultEnv}" env:"ASERTO_ENV" hidden:"" help:"environment override"`
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
