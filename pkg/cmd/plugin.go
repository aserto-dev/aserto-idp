package cmd

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd/plugin"
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
	Config map[string]string //`name:"ceva"`
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
			log.Println(configs[0].Name)

			plugin := Plugin{}
			plugin.Name = idpProvider.GetName()
			plugin.Plugins = append(plugin.Plugins, &PluginFlag{})

			dynamicCommand := kong.DynamicCommand(idpProvider.GetName(), "TODO:This needs to change", "Plugins", &plugin)
			options = append(options, dynamicCommand)
		}
	}
	return options, nil
}
