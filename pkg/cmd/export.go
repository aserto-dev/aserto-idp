package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/directory"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type ExportCmd struct {
	InclUserExt bool
	Source      string
	context     *cc.CC
	kong.Plugins
}

// type HelpStruct struct{}

func (cmd *ExportCmd) Run(c *cc.CC) error {
	cmd.context = c
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewAPIKeyAuth(c.APIKey),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())
	dirClient := conn.DirectoryClient()

	s := make(chan *api.User, 10)
	done := make(chan bool, 1)

	errc := make(chan error, 1)
	go func() {
		for e := range errc {
			c.Log.Debug().Msg(fmt.Sprintf("%s\n", e.Error()))
		}
	}()

	req := &proto.ExportRequest{
		Options: map[string]string{
			"source": cmd.Source,
		},
	}
	go directory.Subscriber(ctx, dirClient, s, done, errc, cmd.InclUserExt)

	client, err := c.Plugin.Export(ctx, req)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	users := []*api.User{}
	for {
		resp, err := client.Recv()
		if err == io.EOF {
			done <- true //means stream is finished
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		log.Printf("Resp received: %s", resp.Data)
		// u := resp.Data.(proto.ExportResponse_User).User
		// users = append(users, &u)
	}

	if users != nil {
		for _, u := range users {
			s <- u
		}
	}

	// close subscriber channel to indicate that the producer done
	close(s)

	// wait for done from subscriber, indicating last received messages has been send
	<-done

	// close error channel as the last action before returning
	close(errc)

	return nil
}
