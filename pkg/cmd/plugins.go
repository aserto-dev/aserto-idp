package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	plugin "github.com/hashicorp/go-plugin"
)

func FindPlugins() map[string]string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dirs = append(dirs, pwd)

	plugins := map[string]string{}
	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, x.PluginPrefix+"*"))
		if err != nil {
			fmt.Println(err)
		}
		for _, f := range files {
			plugins[pluginName(f)] = f
		}
	}

	return plugins
}

func LoadPlugin(name, path string) (grpcplugin.PluginClient, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	// defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return nil, err
	}

	p := raw.(grpcplugin.PluginClient)
	return p, nil
}

func GetPluginHelp(path string) (*shared.HelpMessage, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName(path))
	if err != nil {
		return nil, err
	}

	p := raw.(grpcplugin.PluginClient)
	// p.Info(context)
	// pluginHelp, err := p.Help()
	var res shared.HelpMessage
	fmt.Println(p)
	// structpbconv.Convert(pluginHelp.HelpStruct, &res)
	return &res, nil
}

func pluginName(path string) string {
	file := filepath.Base(path)
	return strings.TrimPrefix(file, x.PluginPrefix)
}
