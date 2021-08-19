package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc/directory"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type ImportCmd struct {
	InclUserExt bool
	Source      string
	kong.Plugins
}

func (cmd *ImportCmd) Run(c *cc.CC) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
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

	go directory.Subscriber(ctx, dirClient, s, done, errc, cmd.InclUserExt)

	users, _ := c.Plugin.LoadUsers(cmd.Source)
	if users != nil {
		for _, u := range users.User {
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
