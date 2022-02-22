package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
	"google.golang.org/protobuf/types/known/structpb"
)

type Plugin struct {
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

var flagsMap = map[string]string{}

func NewPlugin(prov provider.Provider, c *cc.CC) (*Plugin, error) {

	plugin := Plugin{}
	plugin.provider = prov

	pluginClient, err := prov.PluginClient()
	if err != nil {
		return nil, err
	}

	providerInfo, err := pluginClient.Info(c.Context, &proto.InfoRequest{})
	if err != nil {
		return nil, err
	}

	plugin.Name = prov.GetName()

	c.Log.Info().Msgf("loaded plugin %s - version: %s, commit: %s", plugin.Name, providerInfo.Build.Version, providerInfo.Build.Commit)

	for _, config := range providerInfo.Configs {
		configName := fmt.Sprintf("%s-%s", plugin.Name, config.Name)
		flagsMap[configName] = plugin.Name

		plugin.Plugins = append(plugin.Plugins, getFlagStruct(configName, config.Description, plugin.Name, config.Type))
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
				field.Tag = reflect.StructTag(fmt.Sprintf("name:%q help:%q group:%q Flags", flagName, flagDescription, groupName))
			}
		case api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING:
			if field.Type == reflect.TypeOf("string") {
				field.Tag = reflect.StructTag(fmt.Sprintf("name:%q help:%q group:%q Flags", flagName, flagDescription, groupName))
			}
		case api.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER:
			if field.Type == reflect.TypeOf(0) {
				field.Tag = reflect.StructTag(fmt.Sprintf("name:%q help:%q group:%q Flags", flagName, flagDescription, groupName))
			}
		}

		structFields = append(structFields, field)
	}

	newStruct := reflect.StructOf(structFields)

	value := reflect.New(newStruct).Interface()
	return value
}

func getPbStructForNode(pluginName string, pluginConfig map[string]interface{}, node *kong.Node) (*structpb.Struct, error) {
	cliConfigs := getConfigsForNode(pluginName, node)

	for name, value := range pluginConfig {
		if _, ok := cliConfigs[name]; !ok {
			cliConfigs[name] = value
		}
	}

	configStruct, err := structpb.NewStruct(cliConfigs)
	return configStruct, err
}

func getConfigsForNode(pluginName string, node *kong.Node) map[string]interface{} {
	config := make(map[string]interface{})

	for _, flag := range node.Flags {
		if flag.Group == nil || flag.Value.Target.Interface() == flag.DefaultValue.Interface() ||
			!strings.HasPrefix(flag.Name, pluginName) {
			continue
		}

		name := strings.TrimPrefix(flag.Name, pluginName+"-")
		config[name] = flag.Value.Target.Interface()

	}
	return config
}

func validatePlugin(pluginClient grpcplugin.PluginClient, c *cc.CC, config *structpb.Struct, pluginName string, opType proto.OperationType) error {
	c.UI.Note().NoNewline().Msgf("Validating connection to %s", pluginName)

	newFields := make(map[string]*structpb.Value)
	for option, value := range config.Fields {
		newOptionName := strings.TrimPrefix(option, pluginName+"_")
		newFields[newOptionName] = value
	}
	config.Fields = newFields
	validateReq := &proto.ValidateRequest{
		Config: config,
		OpType: opType,
	}

	_, err := pluginClient.Validate(c.Context, validateReq)

	if err == nil {
		c.UI.Success().Msg("Connection validated.")
	}
	return err
}
