package cmd

import (
	"errors"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type GetPluginCmd struct {
	Plugin string `arg:"" short:"n" help:"The idp plugin name and version you want to download. Eg: aserto:x.y.z"`
}

func (cmd *GetPluginCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.Plugin == "" {
		return errors.New("no plugin name was provided")
	}

	info := strings.Split(cmd.Plugin, ":")
	version := "latest"

	if len(info) > 2 {
		return errors.New("plugin is invalid. It must have the following format 'plugin-name:version'")
	}

	if len(info) == 2 {
		version = info[1]
	} else {
		c.Ui.Note().Msg("no version was provided; downloading latest...")
	}

	c.Dispose()
	err := c.Retriever.Download(info[0], version)
	if err != nil {
		return err
	}

	c.Ui.Normal().Msgf("Plugin '%s' '%s' was successfully downloaded", info[0], version)
	return nil
}
