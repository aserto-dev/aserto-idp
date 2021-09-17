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

func GetPluginConfig() []*api.ConfigElement {
	return []*api.ConfigElement{
		{
			Id:          1,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "teststring",
			Description: "Test string config",
			Usage:       "--teststring=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          2,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN,
			Name:        "testbool",
			Description: "Test bool config",
			Usage:       "--testbool",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          3,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER,
			Name:        "testint",
			Description: "Test int config",
			Usage:       "--testint",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
		{
			Id:          4,
			Kind:        api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "teststring2",
			Description: "Test string2 config",
			Usage:       "--teststring2=STRING",
			Mode:        api.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}
}
