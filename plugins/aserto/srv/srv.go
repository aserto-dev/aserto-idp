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
	"github.com/aserto-dev/aserto-idp/plugins/aserto/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type AsertoPluginServer struct{}

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	// response.Build = "placeholder"
	// response.System = ""
	// response.Version = "placeholder"
	response.Description = "Aserto IDP Plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	var dirClient *dir.DirectoryClient

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)
	errc := make(chan error, 1)
	r := make(chan *directory.Result, 1)

	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

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
			if dirClient == nil {
				configBytes, err := protojson.Marshal(req.Config)
				if err != nil {
					errc <- errors.Wrapf(err, "failed to marshal config message")
				}

				config := &config.AsertoConfig{}
				err = json.Unmarshal(configBytes, config)
				if err != nil {
					errc <- errors.Wrapf(err, "failed to unmarshal configs")
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
					errc <- errors.Wrapf(err, "authorizer.Connection")
				}

				ctx = grpcc.SetTenantContext(ctx, tenant)
				dirClient := conn.DirectoryClient()

				go directory.Subscriber(ctx, dirClient, users, r, errc, includeExt)
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
	// Wait for EOF
	<-done

	// close subscriber channel to indicate that the producer done
	close(users)

	// Waif for result from Directory
	result := <-r

	// TODO: (florind) Add more details here
	res := &proto.ImportResponse{}
	if result.Counts != nil {
		res.SuccededCount = result.Counts.Created
		res.FailCount = result.Counts.Errors
	}

	err := srv.SendAndClose(res)
	if err != nil {
		errc <- err
	}

	// close error channel as the last action before returning
	close(errc)

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

	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

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
					User: u,
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

	return nil
}
