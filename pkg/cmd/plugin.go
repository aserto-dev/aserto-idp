package cmd

import (
	"fmt"
	"log"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd/plugin"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/pkg/provider"
	"github.com/aserto-dev/aserto-idp/pkg/provider/finder"
	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type Plugin struct {
	Export plugin.ExportCmd `cmd:""`
	Import plugin.ImportCmd `cmd:""`
	Name   string           `kong:"-"`
	kong.Plugins
}

func (cmd *Plugin) Run(c *cc.CC) error {
	fmt.Println("llll " + cmd.Name + " llll")
	return nil
}

type PluginFlag struct {
	StringFlag string `kong:"-"`
	BoolFlag   bool   `kong:"-"`
	IntFlag    int    `kong:"-"`
}

func SetPluginContext(c *cc.CC, finder ...finder.Finder) error {
	pluginsMap := make(map[string]string)

	for _, finder := range finder {
		pluginPaths, err := finder.Find()
		if err != nil {
			return err
		}

		for _, pluginPath := range pluginPaths {
			idpProvider := provider.NewIDPPluginPlugin(pluginPath)
			pluginName := idpProvider.GetName()

			if path, ok := pluginsMap[pluginName]; ok {
				log.Printf("Plugin %s has already been loaded from %s. Ignoring %s", pluginName, path, pluginPath)
				continue
			}

			pluginsMap[pluginName] = pluginPath

			if pluginName == x.DefaultPluginName {
				client, err := idpProvider.PluginClient()
				if err != nil {
					return err
				}
				c.DefaultIDP = client
			} else if pluginName == c.Command {
				client, err := idpProvider.PluginClient()
				if err != nil {
					return err
				}
				c.CommandIDP = client
			}
		}
	}
	return nil
}

func LoadPlugins(finder ...finder.Finder) ([]kong.Option, error) {
	var options []kong.Option

	pluginsMap := make(map[string]string)

	for _, finder := range finder {
		pluginPaths, err := finder.Find()
		if err != nil {
			return nil, err
		}

		for _, pluginPath := range pluginPaths {
			idpProvider := provider.NewIDPPluginPlugin(pluginPath)
			pluginName := idpProvider.GetName()

			if path, ok := pluginsMap[pluginName]; ok {
				log.Printf("Plugin %s has already been loaded from %s. Ignoring %s", pluginName, path, pluginPath)
				continue
			}

			pluginsMap[pluginName] = pluginPath

			configs, err := idpProvider.Configs()
			if err != nil {
				return nil, err
			}

			plugin := Plugin{}
			for _, config := range configs {
				plugin.Plugins = append(plugin.Plugins, getInterface(config.Name, config.Type))
			}

			dynamicCommand := kong.DynamicCommand(idpProvider.GetName(), "TODO:This needs to change", "Plugins", &plugin)
			options = append(options, dynamicCommand)
		}
	}
	return options, nil
}

func getInterface(flagName string, flagType proto.ConfigElementType) interface{} {
	flag := PluginFlag{}

	flagStructType := reflect.TypeOf(flag)

	structFields := []reflect.StructField{}

	for i := 0; i < flagStructType.NumField(); i++ {
		field := flagStructType.Field(i)

		switch flagType {
		case proto.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN:
			if field.Type == reflect.TypeOf(true) {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" group:"Plugin Flags"`, flagName))
			}
		case proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING:
			if field.Type == reflect.TypeOf("string") {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" group:"Plugin Flags"`, flagName))
			}
		case proto.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER:
			if field.Type == reflect.TypeOf(0) {
				field.Tag = reflect.StructTag(fmt.Sprintf(`name:"%s" group:"Plugin Flags"`, flagName))
			}
		}

		structFields = append(structFields, field)
	}

	newStruct := reflect.StructOf(structFields)

	value := reflect.New(newStruct).Interface()
	return value
}
