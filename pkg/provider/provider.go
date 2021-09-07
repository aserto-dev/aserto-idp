package provider

import (
	"github.com/aserto-dev/aserto-idp/pkg/proto"
)

type Provider interface {
	GetName() string
	GetPath() string
	Configs() ([]*proto.ConfigElement, error)
}
