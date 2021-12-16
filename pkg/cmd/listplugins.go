package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/pkg/errors"
)

type ListPluginsCmd struct {
	Remote bool `short:"r" help:"list plugins available to download from remote"`
}

func (cmd *ListPluginsCmd) Run(context *kong.Context, c *cc.CC) error {
	localPlugins := getLocalPLuginsVersions(c)
	if cmd.Remote {
		idpMajVersion := retriever.IdpMajVersion()
		pluginVersions, err := c.GetRemotePluginsInfo()
		if err != nil {
			return errors.Wrap(err, "failed to retrieve remote information")
		}
		for plugin, majV := range pluginVersions {
			for maj, versions := range majV {
				if maj == idpMajVersion {
					c.Ui.Normal().NoNewline().Msgf("Available versions for '%s'", plugin)
					for _, version := range versions {
						if version == localPlugins[plugin] {
							c.Ui.Normal().NoNewline().Msgf("*\t %s:%s", plugin, version)
						} else {
							c.Ui.Normal().NoNewline().Msgf("\t %s:%s", plugin, version)
						}
					}
					c.Ui.Normal().Msg("")
				}
			}
		}
		return nil
	}

	if len(localPlugins) == 0 {
		c.Ui.Normal().Msg("No local plugins were found")
		return nil
	}
	for name, version := range localPlugins {
		c.Ui.Normal().NoNewline().Msgf("\t %s:%s", name, version)
	}
	c.Ui.Normal().Msg("")

	return nil
}

func getLocalPLuginsVersions(c *cc.CC) map[string]string {
	localVersions := make(map[string]string)
	providers := c.GetProviders()

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
		localVersions[name] = resp.Build.Version
	}

	return localVersions
}
