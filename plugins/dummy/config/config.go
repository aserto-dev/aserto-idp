package config

import "github.com/aserto-dev/aserto-idp/pkg/proto"

func GetPluginConfig() []*proto.ConfigElement {
	return []*proto.ConfigElement{
		{
			Id:          1,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "teststring",
			Description: "Test string config",
			Usage:       "--teststring=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN,
			Name:        "testbool",
			Description: "Test bool config",
			Usage:       "--testbool",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER,
			Name:        "testint",
			Description: "Test int config",
			Usage:       "--testint",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          4,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "teststring2",
			Description: "Test string2 config",
			Usage:       "--teststring2=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
