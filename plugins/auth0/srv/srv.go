package srv

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/auth0/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type Auth0PluginServer struct{}

func (s Auth0PluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	// response.Build = "placeholder"
	// response.System = ""
	// response.Version = "placeholder"
	response.Description = "Auth0 IDP Plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s Auth0PluginServer) Import(srv proto.Plugin_ImportServer) error {
	config := &config.Auth0Config{}
	errc := make(chan error, 1)
	done := make(chan bool, 1)
	subDone := make(chan bool, 1)
	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

	users := make(chan *api.User, 10)
	var sub *Subscriber

	go func() {
		for {
			req, err := srv.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				errc <- errors.Wrapf(err, "srv.Recv()")
			}

			if config.Domain == "" {
				configBytes, err := protojson.Marshal(req.Config)
				if err != nil {
					errc <- err
				}

				err = json.Unmarshal(configBytes, config)
				if err != nil {
					errc <- err
				}
				sub = NewSubscriber(config)
				go sub.Subscriber(users, errc, subDone)
			}

			switch u := req.Data.(type) {
			case *proto.ImportRequest_User:
				{
					users <- u.User
				}
			case *proto.ImportRequest_UserExt:
				{

				}
			}
		}
	}()

	<-done
	close(users)
	<-subDone
	close(errc)
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
