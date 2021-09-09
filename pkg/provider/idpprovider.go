package provider

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	"github.com/hashicorp/go-plugin"
)

type IDPProvider struct {
	Path string
	Name string
}

func NewIDPProvider(path string) Provider {
	asertoPlugin := IDPProvider{
		Path: path,
		Name: providerName(path),
	}
	return &asertoPlugin
}

func (idpProvider *IDPProvider) GetName() string {
	return idpProvider.Name
}

func (idpProvider *IDPProvider) GetPath() string {
	return idpProvider.Path
}

func (idpProvider *IDPProvider) PluginClient() (grpcplugin.PluginClient, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(idpProvider.Path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
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

func (idpProvider *IDPProvider) Info() (*proto.InfoResponse, error) {
	pluginClient, err := idpProvider.PluginClient()
	if err != nil {
		return nil, err
	}

	infoResponse, err := pluginClient.Info(context.Background(), &proto.InfoRequest{})
	if err != nil {
		return nil, err
	}

	return infoResponse, nil
}

func providerName(path string) string {
	file := filepath.Base(path)
	return strings.TrimPrefix(file, x.PluginPrefix)
}
