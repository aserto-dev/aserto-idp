package authorizer

import (
	"context"

	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client gRPC connection
type Client struct {
	conn *grpc.ClientConn
	addr string
}

func Connection(ctx context.Context, addr string, creds credentials.PerRPCCredentials) (*Client, error) {
	gconn, err := grpcc.NewClient(ctx, addr, creds)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &Client{
		conn: gconn.Conn,
		addr: addr,
	}, err
}

// AuthorizerClient -- return authorizer client.
func (c *Client) AuthorizerClient() authz.AuthorizerClient {
	return authz.NewAuthorizerClient(c.conn)
}

// DirectoryClient -- return directory client.
func (c *Client) DirectoryClient() dir.DirectoryClient {
	return dir.NewDirectoryClient(c.conn)
}
