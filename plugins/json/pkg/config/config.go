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

type JsonConfig struct {
	File string `json:"file"`
}

func NewConfig(pbStruct *structpb.Struct) (*JsonConfig, error) {
	configBytes, err := protojson.Marshal(pbStruct)
	if err != nil {
		return nil, err
	}

	config := &JsonConfig{}
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
			Name:        "file",
			Description: "The JSON file",
			Usage:       "--file",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
