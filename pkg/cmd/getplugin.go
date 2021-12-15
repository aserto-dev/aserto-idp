package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
)

type GetPluginCmd struct {
	Plugin string `arg:"" short:"n" help:"The idp plugin name and version you want to download. Eg: aserto:x.y.z"`
}

func (cmd *GetPluginCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.Plugin == "" {
		return errors.New("no plugin name was provided")
	}

	info := strings.Split(cmd.Plugin, ":")

	if len(info) > 2 {
		return errors.New("plugin is invalid. It must have the following format 'plugin-name:version'")
	}

	latest := c.GetLatestVersion(info[0])
	var version string

	if len(info) == 2 {
		version = info[1]
		if version == "latest" {
			version = latest
		}
	} else {
		c.Ui.Note().Msg("no version was provided; downloading latest...")
		version = latest
	}

	if version == "" {
		return fmt.Errorf("couldn't find latest version for %s", info[0])
	}

	provider := c.GetProvider(info[0])
	if provider != nil {
		client, err := provider.PluginClient()
		if err == nil {
			req := &idpplugin.InfoRequest{}
			resp, err := client.Info(c.Context, req)

			if err == nil && resp.Build.Version == version {
				c.Ui.Note().Msgf("Plugin '%s' is already at '%s'", info[0], version)
				return nil
			}
		}
		provider.Kill()
	}

	err := c.Retriever.Download(info[0], version)
	if err != nil {
		return err
	}

	c.Ui.Normal().Msgf("Plugin '%s' '%s' was successfully downloaded", info[0], version)
	return nil
}
