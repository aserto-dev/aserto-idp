package srv

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/directory"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/aserto/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
)

type AsertoPluginServer struct{}

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = "placeholder"
	response.System = ""
	response.Version = "placeholder"
	response.Config = config.GetPluginConfig()

	return &response, nil
}

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	var dirClient *dir.DirectoryClient

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)
	errc := make(chan error, 1)

	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if dirClient == nil {
			authorizerService := req.Options["authorizer"]
			apiKey := req.Options["api_key"]
			tenant := req.Options["tenant"]
			includeExt, err := strconv.ParseBool(req.Options["include_ext"])
			if err != nil {
				return err
			}

			ctx := context.Background()

			conn, err := authorizer.Connection(
				ctx,
				authorizerService,
				grpcc.NewAPIKeyAuth(apiKey),
			)
			if err != nil {
				return err
			}

			ctx = grpcc.SetTenantContext(ctx, tenant)
			dirClient := *conn.DirectoryClient()

			go func() {
				for e := range errc {
					c.Log.Debug().Msg(fmt.Sprintf("%s\n", e.Error()))
				}
			}()

			go directory.Subscriber(ctx, dirClient, users, done, errc, includeExt)
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
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*pluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s AsertoPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	fmt.Println("exporting aserto")
	return nil
}
