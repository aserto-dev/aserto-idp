package plugin

import (
	"github.com/alecthomas/kong"
	"google.golang.org/protobuf/types/known/structpb"
)

func buildStructPb(context *kong.Context) (*structpb.Struct, error) {
	config := make(map[string]interface{})

	flags := context.Selected().Parent.Flags
	for _, flag := range flags {
		config[flag.Name] = flag.Value.Target.Interface()
	}

	configStruct, err := structpb.NewStruct(config)
	return configStruct, err
}
