package provider

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
	sdklogger "github.com/aserto-dev/idp-plugin-sdk/logger"
	sdk "github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog"
)

type IDPProvider struct {
	Path   string
	Name   string
	client *plugin.Client
}

func NewIDPProvider(log *zerolog.Logger, path string) Provider {
	idpProvider := IDPProvider{
		Path: path,
		Name: providerName(path),
	}

	hcpLogger := sdklogger.NewHCLogger(log)

	idpProvider.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         sdk.PluginMap,
		Cmd:             exec.Command(path),
		Logger:          hcpLogger,
		Managed:         false,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	return &idpProvider
}

// Gets the plugin client. Subsequent calls to this will return the same client.
func (idpProvider *IDPProvider) PluginClient() (grpcplugin.PluginClient, error) {

	// Subsequent calls to this will return the same client.
	rpcClient, err := idpProvider.client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin.
	raw, err := rpcClient.Dispense("idp-plugin")
	if err != nil {
		return nil, err
	}

	return raw.(grpcplugin.PluginClient), nil
}

// Kills the plugin process.
func (idpProvider *IDPProvider) Kill() {
	idpProvider.client.Kill()
}

// Gets the name of the provider.
func (idpProvider *IDPProvider) GetName() string {
	return idpProvider.Name
}

// Gets the path of the provider.
func (idpProvider *IDPProvider) GetPath() string {
	return idpProvider.Path
}

func providerName(path string) string {
	file := filepath.Base(path)
	name := strings.TrimPrefix(file, x.PluginPrefix)
	if runtime.GOOS == "windows" {
		name = strings.TrimSuffix(name, ".exe")
	}
	return name
}
