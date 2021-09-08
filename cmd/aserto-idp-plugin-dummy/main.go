package main

import (
	"github.com/aserto-dev/aserto-idp/plugins/dummy/srv"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	plugin "github.com/hashicorp/go-plugin"
)

func main() {
	pSet := make(plugin.PluginSet)
	pSet["idpplugin"] = &grpcplugin.PluginGRPC{
		Impl: &srv.DummyPluginServer{},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pSet,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
