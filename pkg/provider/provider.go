package provider

import (
	"github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
)

type Provider interface {
	GetName() string
	GetPath() string
	PluginClient() (grpcplugin.PluginClient, error)
	Kill()
}
