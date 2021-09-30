package srv

import (
	"context"
	"fmt"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AsertoPluginServer struct{}

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(config.GetVersion),
		Description: "Aserto IDP Plugin",
		Configs:     config.GetPluginConfig(),
	}

	return &response, nil
}

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	var dirClient dir.DirectoryClient

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
				reqConfig := req.GetConfig()
				if reqConfig == nil {
					errc <- status.Error(codes.FailedPrecondition, "Directory service is not initialized")
					continue
				}
				cfg, err := config.NewConfig(reqConfig)
				if err != nil {
					errc <- errors.Wrapf(err, "failed to unmarshal configs")
				}

				ctx := context.Background()

				conn, err := authorizer.Connection(
					ctx,
					cfg.Authorizer,
					grpcc.NewAPIKeyAuth(cfg.ApiKey),
				)
				if err != nil {
					log.Fatalf("Failed to create authorizer connection: %s", err)
				}

				ctx = grpcc.SetTenantContext(ctx, cfg.Tenant)
				dirClient = conn.DirectoryClient()

				go directory.Subscriber(ctx, dirClient, users, r, errc, cfg.IncludeExt)
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

func (s AsertoPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
	return fmt.Errorf("not implemented")
}

// Validate that one use can be retrieved
func (s AsertoPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}

	cfg, err := config.NewConfig(req.Config)
	if err != nil {
		return response, status.Error(codes.InvalidArgument, "failed to parse config")
	}

	conn, err := authorizer.Connection(
		ctx,
		cfg.Authorizer,
		grpcc.NewAPIKeyAuth(cfg.ApiKey),
	)
	if err != nil {
		return response, status.Errorf(codes.Internal, "failed to create authorizar connection %s", err.Error())
	}

	ctx = grpcc.SetTenantContext(ctx, cfg.Tenant)
	dirClient := conn.DirectoryClient()

	_, err = dirClient.ListUsers(ctx, &dir.ListUsersRequest{
		Page: &api.PaginationRequest{
			Size:  1,
			Token: "",
		},
		Base: !cfg.IncludeExt,
	})

	if err != nil {
		return response, status.Errorf(codes.Internal, "failed to get one user: %s", err.Error())
	}
	return response, nil
}

func (s AsertoPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	errc := make(chan error, 1)
	errDone := make(chan bool, 1)

	go grpcerr.SendExportErrors(srv, errc, errDone)

	pageSize := int32(100)
	token := ""

	cfg, err := config.NewConfig(req.Config)
	if err != nil {
		return nil
	}

	ctx := context.Background()

	conn, err := authorizer.Connection(
		ctx,
		cfg.Authorizer,
		grpcc.NewAPIKeyAuth(cfg.ApiKey),
	)
	if err != nil {
		return errors.Wrapf(err, "authorizer.Connection")
	}

	ctx = grpcc.SetTenantContext(ctx, cfg.Tenant)
	dirClient := conn.DirectoryClient()

	for {
		resp, err := dirClient.ListUsers(ctx, &dir.ListUsersRequest{
			Page: &api.PaginationRequest{
				Size:  pageSize,
				Token: token,
			},
			Base: !cfg.IncludeExt,
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
