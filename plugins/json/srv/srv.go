package srv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/json/config"
	"github.com/aserto-dev/aserto-idp/shared/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type JsonPluginServer struct{}

func (s JsonPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = "placeholder"
	response.System = ""
	response.Version = "placeholder"
	response.Description = "Json Plugin"
	response.Configs = config.GetPluginConfig()

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

func (s JsonPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {

	configBytes, err := protojson.Marshal(req.Config)
	if err != nil {
		return err
	}

	config := &config.JsonConfig{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return err
	}

	r, err := os.Open(config.File)
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

func (s JsonPluginServer) Import(srv proto.Plugin_ImportServer) error {
	fmt.Println("import json")
	return nil
}
