package provider

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	"github.com/hashicorp/go-plugin"
)

type IDPProvider struct {
	Path   string
	Name   string
	client *plugin.Client
}

func NewIDPProvider(path string) Provider {
	idpProvider := IDPProvider{
		Path: path,
		Name: providerName(path),
	}

	idpProvider.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	return &idpProvider
}

// Gets the plugin client. Subsequent calls to this will return the same client
func (idpProvider *IDPProvider) PluginClient() (grpcplugin.PluginClient, error) {

	// Subsequent calls to this will return the same client.
	rpcClient, err := idpProvider.client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("idp-plugin")
	if err != nil {
		return nil, err
	}

	return raw.(grpcplugin.PluginClient), nil
}

// Kills the plugin process
func (idpProvider *IDPProvider) Kill() {
	idpProvider.client.Kill()
}

// Gets the name of the provider
func (idpProvider *IDPProvider) GetName() string {
	return idpProvider.Name
}

// Gets the path of the provider
func (idpProvider *IDPProvider) GetPath() string {
	return idpProvider.Path
}

func providerName(path string) string {
	file := filepath.Base(path)
	return strings.TrimPrefix(file, x.PluginPrefix)
}
