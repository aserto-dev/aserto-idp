package plugin

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	"google.golang.org/protobuf/types/known/structpb"
)

func getPbStructForNode(pluginConfig map[string]interface{}, node *kong.Node) (*structpb.Struct, error) {
	cliConfigs := getConfigsForNode(node)

	for name, value := range pluginConfig {
		if _, ok := cliConfigs[name]; !ok {
			cliConfigs[name] = value
		}
	}

	configStruct, err := structpb.NewStruct(cliConfigs)
	return configStruct, err
}

func getConfigsForNode(node *kong.Node) map[string]interface{} {
	config := make(map[string]interface{})

	for _, flag := range node.Flags {
		// CLI flags do not have groups
		if flag.Group != nil && flag.Value.Target.Interface() != flag.DefaultValue.Interface() {
			config[flag.Name] = flag.Value.Target.Interface()
		}
	}
	return config
}

func validatePlugin(pluginClient grpcplugin.PluginClient, c *cc.CC, config *structpb.Struct, pluginName string) error {
	c.Ui.Note().NoNewline().Msgf("Validating connection to %s", pluginName)
	validateReq := &proto.ValidateRequest{
		Config: config,
	}

	_, err := pluginClient.Validate(c.Context, validateReq)

	if err == nil {
		c.Ui.Success().Msg("Connection validated.")
	}
	return err
}
