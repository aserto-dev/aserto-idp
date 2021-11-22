package cmd

import (
	"errors"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type GetPluginCmd struct {
	Name    string `short:"n" help:"The idp plugin name you want to download"`
	Version string `short:"n" help:"The idp plugin version you want to download"`
}

func (cmd *GetPluginCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.Name == "" {
		return errors.New("no plugin name was provided")
	}

	if cmd.Version == "" {
		cmd.Version = "latest"
		c.Log.Warn().Msg("no version was provided; downloading latest...")
	}

	ghcr := c.Retriever
	err := ghcr.Connect()
	if err != nil {
		return err
	}
	err = ghcr.Download(cmd.Name, cmd.Version)
	if err != nil {
		return err
	}
	return nil
}
