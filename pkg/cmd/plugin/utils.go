package plugin

import (
	"github.com/alecthomas/kong"
	"google.golang.org/protobuf/types/known/structpb"
)

func getPbStructForNode(node *kong.Node) (*structpb.Struct, error) {
	config := make(map[string]interface{})

	for _, flag := range node.Flags {
		// CLI flags do not have groups
		if flag.Group != nil {
			config[flag.Name] = flag.Value.Target.Interface()
		}
	}
	configStruct, err := structpb.NewStruct(config)
	return configStruct, err
}
