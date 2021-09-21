package srv

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/dummy/pkg/config"
	"github.com/aserto-dev/aserto-idp/shared/version"
)

type DummyPluginServer struct{}

func (s DummyPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(config.GetVersion),
		Description: "Dummy plugin",
		Configs:     config.GetPluginConfig(),
	}

	return &response, nil
}

func (s DummyPluginServer) Import(srv proto.Plugin_ImportServer) error {
	log.Println("not implemented")
	return nil
}

// func (s DummyPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

func (DummyPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}
	return response, nil
}

func (s DummyPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	log.Println("not implemented")
	return nil
}
