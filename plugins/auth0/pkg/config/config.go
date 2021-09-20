package config

import api "github.com/aserto-dev/go-grpc/aserto/api/v1"

// values set by linker using ldflag -X
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

type Auth0Config struct {
	Domain       string `json:"domain"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func GetPluginConfig() []*api.ConfigElement {
	return []*api.ConfigElement{
		{
			Id:          1,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "domain",
			Description: "Auth0 domain",
			Usage:       "--domain=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "client_id",
			Description: "The Client ID",
			Usage:       "--client_id=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "client_secret",
			Description: "The Client Secret",
			Usage:       "--client_secret=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}