package cmd

import (
	"errors"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type GetPluginCmd struct {
	Name string `short:"n" help:"The idp plugin name and version you want to download. Eg: aserto:x.y.z"`
}

func (cmd *GetPluginCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.Name == "" {
		return errors.New("no plugin name was provided")
	}

	info := strings.Split(cmd.Name, ":")
	version := "latest"

	if len(info) > 2 {
		return errors.New("please provide plugin to download as: plugin-name:version")
	}

	if len(info) == 2 {
		version = info[1]
	} else {
		c.Ui.Exclamation().Msg("no version was provided; downloading latest...")
	}

	c.Dispose()
	err := c.Retriever.Download(info[0], version)
	if err != nil {
		return err
	}

	c.Ui.Normal().Msgf("Plugin %s %s was successfully downloaded", info[0], version)
	return nil
}
