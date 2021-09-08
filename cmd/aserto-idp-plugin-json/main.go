package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	pSet["idp-plugin"] = &grpcplugin.PluginGRPC{
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

func (s PluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = "placeholder"
	response.System = ""
	response.Version = "placeholder"
	response.Config = []*proto.ConfigElement{
		{
			Id:          1,
			Kind:        proto.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE,
			Type:        proto.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING,
			Name:        "file",
			Description: "The JSON file",
			Usage:       "--file",
			Mode:        proto.DisplayMode_DISPLAY_MODE_NORMAL,
			ReadOnly:    false,
		},
	}

	return &response, nil
}

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
		if err = srv.Send(res); err != nil {
			log.Println(err)
			return err
		}
	}

	if _, err = dec.Token(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s PluginServer) Import(srv proto.Plugin_ImportServer) error {
	fmt.Println("import json")
	return nil
}
