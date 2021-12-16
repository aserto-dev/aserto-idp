package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
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
	pluginsVersions, err := c.GetRemotePluginsInfo()
	if err != nil {
		return errors.New("failed to get remote information")
	}
	if pluginsVersions[pluginName] == nil {
		return fmt.Errorf("plugin '%s' does not exists", pluginName)
	}

	latestVersion, err := c.GetLatestVersion(pluginName)
	if err != nil {
		return errors.Wrapf(err, "failed to get versions for '%s'", pluginName)
	}
	if latestVersion == "" {
		return fmt.Errorf("couldn't find latest version for '%s'", pluginName)
	}

	err = c.Retriever.Download(pluginName, latestVersion)
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

func checkForUpdates(provider provider.Provider, c *cc.CC) (bool, string, error) {
	client, err := provider.PluginClient()
	if err != nil {
		return false, "", errors.Wrap(err, "failed to get plugin client")
	}
	req := &idpplugin.InfoRequest{}
	resp, err := client.Info(c.Context, req)
	if err != nil {
		return false, "", errors.Wrap(err, "failed to get plugin info")
	}

	latest, err := c.GetLatestVersion(provider.GetName())
	if err != nil {
		return false, "", errors.Wrap(err, "failed to get remote information about plugins")
	}
	if latest == "" {
		return false, "", err
	}

	currentVers := strings.Split(resp.Build.Version, ".")
	latestVers := strings.Split(latest, ".")

	for index, ver := range currentVers {
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
