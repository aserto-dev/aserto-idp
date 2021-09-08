package config

import "github.com/aserto-dev/aserto-idp/pkg/proto"

func GetPluginConfig() []*proto.ConfigElement {
	return []*proto.ConfigElement{
		{
			Id:          1,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "file",
			Description: "The JSON file",
			Usage:       "--file",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
