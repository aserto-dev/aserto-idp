// Package shared contains shared data between the host and plugins.
package shared

import (
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	plugin "github.com/hashicorp/go-plugin"
)

type HelpMessage struct {
	Json string
}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = plugin.PluginSet{
	"idp-plugin": &grpcplugin.PluginGRPC{},
}

// Provider is the interface that we're exposing as a plugin.
// type Provider interface {
// 	LoadUsers(source string) (*proto.LoadUsersResponse, error)
// 	Help() (*proto.HelpResponse, error)
// }

// // This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
// type ProviderPlugin struct {
// 	// GRPCPlugin must still implement the Plugin interface
// 	plugin.Plugin
// 	// Concrete implementation, written in Go. This is only used for plugins
// 	// that are written in Go.
// 	Impl Provider
// }

// func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
// 	// proto.RegisterProviderServer(s, &GRPCServer{Impl: p.Impl})
// 	proto.RegisterProviderServer(s, &GRPCServer{Impl: p.Impl})
// 	return nil
// }

// func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
// 	return &GRPCClient{client: proto.NewProviderClient(c)}, nil
// }
