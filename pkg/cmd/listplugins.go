package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
)

type ListPluginsCmd struct {
	Remote bool `short:"r" help:"list plugins available to download from remote"`
}

func (cmd *ListPluginsCmd) Run(context *kong.Context, c *cc.CC) error {
	if cmd.Remote {
		idpMajVersion := retriever.IdpMajVersion()
		pluginVersions := retriever.PluginVersions(c.Retriever)
		for plugin, majV := range pluginVersions {
			for maj, versions := range majV {
				if maj == idpMajVersion {
					c.Ui.Normal().Msgf("Plugin %s", plugin)
					for _, version := range versions {
						c.Ui.Normal().Msgf("%s", version)
					}
				}
			}
		}
	} else {
		providers := c.GetProviders()
		if len(providers) == 0 {
			c.Ui.Normal().Msg("No local plugins were found")
			return nil
		}
		for name, provider := range providers {
			client, err := provider.PluginClient()
			if err != nil {
				c.Log.Debug().Msgf("failed to get client for plugin '%s'", name)
				continue
			}
			req := &idpplugin.InfoRequest{}
			resp, err := client.Info(c.Context, req)
			if err != nil {
				c.Log.Debug().Msgf("failed to retrieve info on plugin '%s'", name)
				continue
			}
			c.Ui.Normal().Msgf("Plugin %s", name)
			c.Ui.Normal().Msgf("%s", resp.Build.Version)
		}
	}
	return nil
}
