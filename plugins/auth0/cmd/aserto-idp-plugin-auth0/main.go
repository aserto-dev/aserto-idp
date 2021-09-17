package main

import (
	"github.com/aserto-dev/aserto-idp/plugins/auth0/pkg/srv"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	pSet := make(plugin.PluginSet)
	pSet["idp-plugin"] = &grpcplugin.PluginGRPC{
		Impl: &srv.Auth0PluginServer{},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pSet,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
