package plugin

import (
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
)

type ImportCmd struct {
}

func (cmd *ImportCmd) Run(app *kong.Kong, context *kong.Context, c *cc.CC) error {

	defaultProviderName := c.GetDefaultProvider().GetName()
	defaultProviderConfigs, err := getPbStructForNode(c.Config.Plugins[defaultProviderName], context.Path[0].Node())
	if err != nil {
		return err
	}

	//TODO: Handle this
	includeExt := false

	providerName := context.Selected().Parent.Name
	providerConfigs, err := getPbStructForNode(c.Config.Plugins[providerName], context.Selected().Parent)
	if err != nil {
		return err
	}
	req := &proto.ExportRequest{
		Config: providerConfigs,
	}

	providerClient, err := c.GetProvider(providerName).PluginClient()
	if err != nil {
		return err
	}
	exportClient, err := providerClient.Export(c.Context, req)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	defaultProviderClient, err := c.GetDefaultProvider().PluginClient()
	if err != nil {
		return err
	}

	importClient, err := defaultProviderClient.Import(c.Context)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	users := make(chan *api.User, 10)
	done := make(chan bool, 1)
	errc := make(chan error, 1)
	result := make(chan *proto.ImportResponse, 1)

	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

	// send users
	go func() {
		for user := range users {
			if !includeExt {
				user.Attributes = &api.AttrSet{}
				user.Applications = make(map[string]*api.AttrSet)
			}

			req := &proto.ImportRequest{
				Config: defaultProviderConfigs,
				Data: &proto.ImportRequest_User{
					User: user,
				},
			}
			if err = importClient.Send(req); err != nil {
				errc <- errors.Wrapf(err, "stream send %s", user.Id)
			}
		}
		res, err := importClient.CloseAndRecv()
		if err != nil {
			errc <- errors.Wrapf(err, "stream.CloseAndRecv()")
		}
		result <- res
	}()

	// receive users
	go func() {
		for {
			resp, err := exportClient.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("Resp received: %s", resp.Data)
			switch u := resp.Data.(type) {
			case *proto.ExportResponse_User:
				{
					users <- u.User
				}
			case *proto.ExportResponse_UserExt:
				{

				}
			}
		}
	}()

	// Wait for EOF
	<-done

	close(users)

	// Wait for Result from import
	res := <-result

	close(errc)

	if res != nil {
		log.Printf("Succeeded: %d\n", res.SuccededCount)
		log.Printf("Failed: %d\n", res.FailCount)
	}
	return nil
}
