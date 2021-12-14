package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/aserto-dev/aserto-idp/pkg/version"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/pkg/errors"
)

type CLI struct {
	Config      string         `short:"c" type:"path" help:"Path to the config file. Any argument provided to the CLI will take precedence."`
	Delete      DeleteCmd      `cmd:"" help:"delete user ids from an user-provider idp"`
	Exec        ExecCmd        `cmd:"" help:"import users from an user-provided idp to another user-provided idp"`
	GetPlugin   GetPluginCmd   `cmd:"" help:"download plugin"`
	ListPlugins ListPluginsCmd `cmd:"" help:"list available plugins"`
	Version     VersionCmd     `cmd:"" help:"version information"`
	Verbosity   int            `short:"v" type:"counter" help:"Use to increase output verbosity."`
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

func downloadProvider(pluginName string, c *cc.CC) error {
	pluginsVersions := retriever.PluginVersions(c.Retriever)
	if pluginsVersions[pluginName] == nil {
		return fmt.Errorf("unknown plugin name %s was provided", pluginName)
	}
	err := c.Retriever.Download(pluginName, "latest")
	if err != nil {
		return err
	}
	err = c.LoadProviders()
	if err != nil {
		return err
	}
	if !c.ProviderExists(pluginName) {
		return fmt.Errorf("plugin %s couldn't be downloaded", pluginName)
	}
	return nil
}

func checkForUpdates(provider provider.Provider, store retriever.Retriever) (bool, string, error) {
	client, err := provider.PluginClient()
	if err != nil {
		return false, "", errors.Wrap(err, "can't get client")
	}
	req := &idpplugin.InfoRequest{}
	resp, err := client.Info(context.Background(), req)
	if err != nil {
		return false, "", errors.Wrap(err, "can't get version")
	}

	pluginsVersions := retriever.PluginVersions(store)
	availableVersions :=
		pluginsVersions[provider.GetName()][retriever.IdpMajVersion()]

	presentVers := strings.Split(resp.Build.Version, ".")
	latestVers := strings.Split(availableVersions[0], ".")

	for index, ver := range presentVers {
		intVer, err := strconv.Atoi(ver)
		if err != nil {
			return false, "", err
		}
		intLatestPart, err := strconv.Atoi(latestVers[index])
		if err != nil {
			return false, "", err
		}
		if intLatestPart > intVer {
			return true, strings.Join(latestVers, "."), nil
		}
	}
	return false, "", nil
}
