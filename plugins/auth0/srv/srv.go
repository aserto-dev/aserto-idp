package srv

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/auth0/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type Auth0PluginServer struct{}

func (s Auth0PluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = "placeholder"
	response.System = ""
	response.Version = "placeholder"
	response.Description = "Auth0 IDP Plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s Auth0PluginServer) Import(srv proto.Plugin_ImportServer) error {
	return nil
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*pluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s Auth0PluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	configBytes, err := protojson.Marshal(req.Config)
	if err != nil {
		return err
	}

	config := &config.Auth0Config{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return err
	}
	log.Println(config)
	errc := make(chan error, 1)
	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)

	go func() {
		for u := range users {
			res := &proto.ExportResponse{
				Data: &proto.ExportResponse_User{
					User: u,
				},
			}
			if err = srv.Send(res); err != nil {
				errc <- err
			}
		}
		done <- true
	}()

	p := NewProducer(config)
	p.Producer(users, errc)

	close(users)

	<-done

	close(errc)

	return nil
}
