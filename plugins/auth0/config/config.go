package config

import "github.com/aserto-dev/aserto-idp/pkg/proto"

type Auth0Config struct {
	Domain       string `json:"domain"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func GetPluginConfig() []*proto.ConfigElement {
	return []*proto.ConfigElement{
		{
			Id:          1,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "domain",
			Description: "Auth0 domain",
			Usage:       "--domain=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "client_id",
			Description: "The Client ID",
			Usage:       "--client_id=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "client_secret",
			Description: "The Client Secret",
			Usage:       "--client_secret=STRING",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
