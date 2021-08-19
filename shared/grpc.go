package shared

import (
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	client proto.ProviderClient
}

func (m *GRPCClient) LoadUsers(source string) (*proto.LoadUsersResponse, error) {
	res, err := m.client.LoadUsers(context.Background(), &proto.LoadUsersRequest{
		Source: source,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *GRPCClient) Help() (*proto.HelpResponse, error) {
	resp, err := m.client.Help(context.Background(), &proto.HelpRequest{})
	return resp, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Provider
}

func (m *GRPCServer) Help(ctx context.Context, req *proto.HelpRequest) (*proto.HelpResponse, error) {
	resp, err := m.Impl.Help()
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *GRPCServer) LoadUsers(ctx context.Context, req *proto.LoadUsersRequest) (*proto.LoadUsersResponse, error) {
	res, err := m.Impl.LoadUsers(req.Source)
	if err != nil {
		return nil, err
	}
	return res, nil
}
