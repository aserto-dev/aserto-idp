package main

import (
	"github.com/aserto-dev/aserto-idp/plugins/auth0/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
	sdk "github.com/aserto-dev/idp-plugin-sdk/plugin"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	pSet := make(plugin.PluginSet)
	pSet["idp-plugin"] = &grpcplugin.PluginGRPC{
		Impl: &srv.Auth0PluginServer{},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.Handshake,
		Plugins:         pSet,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
