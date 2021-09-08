package provider

import (
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
)

type Provider interface {
	GetName() string
	GetPath() string
	Configs() ([]*proto.ConfigElement, error)
	PluginClient() (grpcplugin.PluginClient, error)
}
