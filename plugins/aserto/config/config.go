package config

import "github.com/aserto-dev/aserto-idp/pkg/proto"

func GetPluginConfig() []*proto.ConfigElement {
	return []*proto.ConfigElement{
		{
			Id:          1,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "authorizer",
			Description: "The Authorizer endpoint",
			Usage:       "--authorizer=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "tenant",
			Description: "The Tenant ID",
			Usage:       "--tenant=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "api_key",
			Description: "Aserto API Key",
			Usage:       "--api_key=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          4,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN,
			Name:        "include_ext",
			Description: "Include User Extensions",
			Usage:       "--include_ext=BOOL",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
