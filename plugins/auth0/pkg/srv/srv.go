package srv

import (
	"context"
	"io"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/auth0/pkg/config"
	grpcerr "github.com/aserto-dev/aserto-idp/shared/errors"
	"github.com/aserto-dev/aserto-idp/shared/version"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/auth0.v5/management"
)

type Auth0PluginServer struct{}

func (s Auth0PluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(config.GetVersion),
		Description: "Auth0 IDP Plugin",
		Configs:     config.GetPluginConfig(),
	}

	return &response, nil
}

func (s Auth0PluginServer) Import(srv proto.Plugin_ImportServer) error {
	var cfg *config.Auth0Config
	errc := make(chan error, 1)
	done := make(chan bool, 1)
	subDone := make(chan bool, 1)
	errDone := make(chan bool, 1)

	go grpcerr.SendImportErrors(srv, errc, errDone)

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

			if cfg == nil {
				cfg, err = config.NewConfig(req.GetConfig())
				if err != nil {
					errc <- err
				}
				sub = NewSubscriber(cfg)
				go sub.Subscriber(users, errc, subDone)
			}

			if user := req.GetUser(); user != nil {
				if u := user.GetUser(); u != nil {
					users <- u
				}
			}
		}
	}()

	<-done
	close(users)
	<-subDone
	close(errc)
	<-errDone
	return nil
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

func (Auth0PluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}
	cfg, err := config.NewConfig(req.Config)
	if err != nil {
		return response, status.Error(codes.InvalidArgument, "failed to parse config")
	}

	mgnt, err := management.New(
		cfg.Domain,
		management.WithClientCredentials(
			cfg.ClientID,
			cfg.ClientSecret,
		))
	if err != nil {
		return response, status.Error(codes.Internal, "failed to connect to Auth0")
	}

	_, err = mgnt.Connection.ReadByName("Username-Password-Authentication")
	if err != nil {
		return response, status.Error(codes.Internal, "failed to get Auth0 connection")
	}

	return response, nil
}

func (s Auth0PluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {

	cfg, err := config.NewConfig(req.Config)
	if err != nil {
		return err
	}

	errc := make(chan error, 1)
	errDone := make(chan bool, 1)

	go grpcerr.SendExportErrors(srv, errc, errDone)

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)

	go func() {
		for u := range users {
			res := &proto.ExportResponse{
				Data: &proto.ExportResponse_User{
					User: &proto.User{
						Data: &proto.User_User{
							User: u,
						},
					},
				},
			}
			if err = srv.Send(res); err != nil {
				errc <- err
			}
		}
		done <- true
	}()

	p := NewProducer(cfg)
	p.Producer(users, errc)

	close(users)

	<-done

	close(errc)
	<-errDone

	return nil
}
