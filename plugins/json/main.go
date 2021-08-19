package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/protobuf/types/known/structpb"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type Provider struct{}

// var help struct {
// 	Json string
// }

var help = &structpb.Struct{
	Fields: map[string]*structpb.Value{
		"Json": &structpb.Value{
			Kind: &structpb.Value_StringValue{},
		},
	},
}

func (Provider) LoadUsers(source string) (*proto.LoadUsersResponse, error) {
	log.Println(source)
	r, err := os.Open(source)
	if err != nil {
		// errc <- errors.Wrapf(err, "open %s", source)
	}

	dec := json.NewDecoder(r)

	if _, err = dec.Token(); err != nil {
		// errc <- errors.Wrapf(err, "token open")
	}

	users := []*api.User{}
	// u := api.User{
	// 	Id:    "asd",
	// 	Email: "agfsd",
	// }
	// users = append(users, &u)
	for dec.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(dec, &u); err != nil {
			// errc <- errors.Wrapf(err, "unmarshal next")
		}
		users = append(users, &u)
		// s <- &u
		// p.count++
	}

	if _, err = dec.Token(); err != nil {
		// errc <- errors.Wrapf(err, "token close")
	}
	res := &proto.LoadUsersResponse{
		User: users,
	}
	return res, nil
}

func (Provider) Help() (*proto.HelpResponse, error) {
	resp := &proto.HelpResponse{
		HelpStruct: help,
	}

	return resp, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"json": &shared.ProviderPlugin{Impl: &Provider{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
