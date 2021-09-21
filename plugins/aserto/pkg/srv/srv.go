package srv

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/directory"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/aserto/pkg/config"
	grpcerr "github.com/aserto-dev/aserto-idp/shared/errors"
	"github.com/aserto-dev/aserto-idp/shared/version"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type AsertoPluginServer struct{}

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = version.GetBuildInfo(config.GetVersion)
	response.Description = "Aserto IDP Plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	var dirClient *dir.DirectoryClient

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)
	errDone := make(chan bool, 1)
	errc := make(chan error, 1)
	r := make(chan *directory.Result, 1)

	go grpcerr.SendImportErrors(srv, errc, errDone)

	go func() {
		for {
			req, err := srv.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Println(errors.Wrapf(err, "srv.Recv()"))
			}
			if dirClient == nil {
				if cfg := req.GetConfig(); cfg != nil {
					configBytes, err := protojson.Marshal(cfg)
					if err != nil {
						log.Println(errors.Wrapf(err, "failed to marshal config message"))
					}

					config := &config.AsertoConfig{}
					err = json.Unmarshal(configBytes, config)
					if err != nil {
						log.Println(errors.Wrapf(err, "failed to unmarshal configs"))
					}

					authorizerService := config.Authorizer
					apiKey := config.ApiKey
					tenant := config.Tenant
					includeExt := config.IncludeExt

					ctx := context.Background()

					conn, err := authorizer.Connection(
						ctx,
						authorizerService,
						grpcc.NewAPIKeyAuth(apiKey),
					)
					if err != nil {
						log.Println(errors.Wrapf(err, "authorizer.Connection"))
					}

					ctx = grpcc.SetTenantContext(ctx, tenant)
					dirClient := conn.DirectoryClient()

					go directory.Subscriber(ctx, dirClient, users, r, errc, includeExt)
				}
			}

			if user := req.GetUser(); user != nil {
				if u := user.GetUser(); u != nil {
					users <- u
				}
			}
		}
	}()
	// Wait for EOF
	<-done

	// close subscriber channel to indicate that the producer done
	close(users)

	// Waif for result from Directory
	result := <-r

	// TODO: (florind) Add more details here
	if result.Counts != nil {
		log.Printf("Created: %d\n", result.Counts.Created)
		log.Printf("Failed: %d\n", result.Counts.Errors)
	}

	// close error channel as the last action before returning
	close(errc)
	<-errDone

	return result.Err
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*pluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s AsertoPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	errc := make(chan error, 1)
	errDone := make(chan bool, 1)

	go grpcerr.SendExportErrors(srv, errc, errDone)

	pageSize := int32(100)
	token := ""

	configBytes, err := protojson.Marshal(req.Config)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal config message")
	}

	config := &config.AsertoConfig{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal configs")
	}

	authorizerService := config.Authorizer
	apiKey := config.ApiKey
	tenant := config.Tenant
	includeExt := config.IncludeExt

	ctx := context.Background()

	conn, err := authorizer.Connection(
		ctx,
		authorizerService,
		grpcc.NewAPIKeyAuth(apiKey),
	)
	if err != nil {
		return errors.Wrapf(err, "authorizer.Connection")
	}

	ctx = grpcc.SetTenantContext(ctx, tenant)
	dirClient := conn.DirectoryClient()

	for {
		resp, err := dirClient.ListUsers(ctx, &dir.ListUsersRequest{
			Page: &api.PaginationRequest{
				Size:  pageSize,
				Token: token,
			},
			Base: !includeExt,
		})

		if err != nil {
			return errors.Wrapf(err, "list users")
		}
		for _, u := range resp.Results {
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

		if resp.Page.NextToken == "" {
			break
		}

		token = resp.Page.NextToken
	}

	close(errc)
	<-errDone

	return nil
}
