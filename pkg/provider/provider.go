package provider

import (
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
)

type Provider interface {
	GetName() string
	GetPath() string
	PluginClient() (grpcplugin.PluginClient, error)
	Kill()
}
