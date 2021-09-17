package srv

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/dummy/config"
	"github.com/aserto-dev/aserto-idp/shared/version"
)

type DummyPluginServer struct{}

func (s DummyPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = version.GetBuildInfo(config.GetVersion)
	response.Description = "Dummy plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s DummyPluginServer) Import(srv proto.Plugin_ImportServer) error {
	log.Println("not implemented")
	return nil
}

// func (s DummyPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*DummyPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s DummyPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	log.Println("not implemented")
	return nil
}
