package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	"github.com/aserto-dev/aserto-idp/shared/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	plugin "github.com/hashicorp/go-plugin"
)

// type Provider struct{}

// var help = &structpb.Struct{
// 	Fields: map[string]*structpb.Value{
// 		"Json": &structpb.Value{
// 			Kind: &structpb.Value_StringValue{},
// 		},
// 	},
// }

// func (Provider) LoadUsers(source string) (*proto.LoadUsersResponse, error) {
// 	log.Println(source)
// 	r, err := os.Open(source)
// 	if err != nil {
// 		log.Println(err)
// 		// errc <- errors.Wrapf(err, "open %s", source)
// 	}

// 	dec := json.NewDecoder(r)

// 	if _, err = dec.Token(); err != nil {
// 		log.Println(err)
// 		// errc <- errors.Wrapf(err, "token open")
// 	}

// 	users := []*api.User{}
// 	for dec.More() {
// 		u := api.User{}
// 		if err := pb.UnmarshalNext(dec, &u); err != nil {
// 			log.Println(err)
// 			// errc <- errors.Wrapf(err, "unmarshal next")
// 		}
// 		users = append(users, &u)
// 	}

// 	if _, err = dec.Token(); err != nil {
// 		log.Println(err)
// 		// errc <- errors.Wrapf(err, "token close")
// 	}
// 	res := &proto.LoadUsersResponse{
// 		User: users,
// 	}
// 	return res, nil
// }

// func (Provider) Help() (*proto.HelpResponse, error) {
// 	resp := &proto.HelpResponse{
// 		HelpStruct: help,
// 	}

// 	return resp, nil
// }

func main() {
	pSet := make(plugin.PluginSet)
	pSet["json"] = &grpcplugin.PluginGRPC{
		Impl: &PluginServer{},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pSet,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type PluginServer struct{}

// func (s pluginServer) Help(ctx context.Context, req *proto.HelpRequest) (*proto.HelpResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

// func (s pluginServer) Import(srv proto.Plugin_ImportServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*pluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s PluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	r, err := os.Open(req.Options["source"])
	if err != nil {
		log.Println(err)
		return err
	}

	dec := json.NewDecoder(r)

	if _, err = dec.Token(); err != nil {
		log.Println(err)
		return err
	}

	for dec.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(dec, &u); err != nil {
			log.Println(err)
			return err
		}
		res := &proto.ExportResponse{
			Data: &proto.ExportResponse_User{
				User: &u,
			},
		}
		srv.Send(res)
	}

	if _, err = dec.Token(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
