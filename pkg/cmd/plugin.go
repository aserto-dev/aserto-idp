package cmd

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd/plugin"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
)

type Plugin struct {
	Export      plugin.ExportCmd  `cmd:""`
	Import      plugin.ImportCmd  `cmd:""`
	provider    provider.Provider `kong:"-"`
	Description string            `kong:"-"`
	Name        string            `kong:"-"`
	kong.Plugins
}

func (plugin *Plugin) Run(c *cc.CC) error {
	return nil
}

type PluginFlag struct {
	StringFlag string `kong:"-"`
	BoolFlag   bool   `kong:"-"`
	IntFlag    int    `kong:"-"`
}

func NewPlugin(provider provider.Provider, c *cc.CC) (*Plugin, error) {

	plugin := Plugin{}
	plugin.provider = provider

	pluginClient, err := provider.PluginClient()
	if err != nil {
		return nil, err
	}

	providerInfo, err := pluginClient.Info(c.Context, &proto.InfoRequest{})
	if err != nil {
		return nil, err
	}

	plugin.Name = provider.GetName()

	c.Log.Info().Msgf("loaded plugin %s - version: %s, commit: %s", plugin.Name, providerInfo.Build.Version, providerInfo.Build.Commit)

	for _, config := range providerInfo.Configs {
		plugin.Plugins = append(plugin.Plugins, getFlagStruct(config.Name, config.Description, plugin.Name, config.Type))
	}

	plugin.Description = providerInfo.Description

	return &plugin, nil
}

func getFlagStruct(flagName, flagDescription, groupName string, flagType api.ConfigElementType) interface{} {
	flag := PluginFlag{}

	flagStructType := reflect.TypeOf(flag)

	structFields := []reflect.StructField{}

	for i := 0; i < flagStructType.NumField(); i++ {
		field := flagStructType.Field(i)

		switch flagType {
		case api.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN:
			if field.Type == reflect.TypeOf(true) {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" help:"%s" group:"%s Flags"`, flagName, flagDescription, groupName))
			}
		case api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING:
			if field.Type == reflect.TypeOf("string") {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" help:"%s" group:"%s Flags"`, flagName, flagDescription, groupName))
			}
		case api.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER:
			if field.Type == reflect.TypeOf(0) {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" help:"%s" group:"%s Flags"`, flagName, flagDescription, groupName))
			}
		}

		structFields = append(structFields, field)
	}

	newStruct := reflect.StructOf(structFields)

	value := reflect.New(newStruct).Interface()
	return value
}
