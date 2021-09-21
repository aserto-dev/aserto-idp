package config

import (
	"encoding/json"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// values set by linker using ldflag -X
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

type AsertoConfig struct {
	Authorizer string `json:"authorizer"`
	Tenant     string `json:"tenant"`
	ApiKey     string `json:"api_key"`
	IncludeExt bool   `json:"include_ext"`
}

func NewConfig(pbStruct *structpb.Struct) (*AsertoConfig, error) {
	configBytes, err := protojson.Marshal(pbStruct)
	if err != nil {
		return nil, err
	}

	config := &AsertoConfig{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetPluginConfig() []*api.ConfigElement {
	return []*api.ConfigElement{
		{
			Id:          1,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "authorizer",
			Description: "The Authorizer endpoint",
			Usage:       "--authorizer=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "tenant",
			Description: "The Tenant ID",
			Usage:       "--tenant=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "api_key",
			Description: "Aserto API Key",
			Usage:       "--api_key=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          4,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN,
			Name:        "include_ext",
			Description: "Include User Extensions",
			Usage:       "--include_ext=BOOL",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
