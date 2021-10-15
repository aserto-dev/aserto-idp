package cmd

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteCmd struct {
	From    string   `short:"f" help:"The idp provider name you want to delete from"`
	UserIds []string `arg:"" name:"user_id" help:"Users to remove." type:"string"`
}

func (cmd *DeleteCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.From == "" || !c.ProviderExists(cmd.From) {
		return status.Error(codes.InvalidArgument, "no \"--from\" idp or an unavailable idp was provided")
	}

	providerName := cmd.From
	providerConfigs, err := getPbStructForNode(c.Config.Plugins[providerName], context.Path[0].Node())
	if err != nil {
		return err
	}

	providerClient, err := c.GetProvider(providerName).PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(providerClient, c, providerConfigs, providerName)
	if err != nil {
		return err
	}

	deleteClient, err := providerClient.Delete(c.Context)
	if err != nil {
		return err
	}

	doneDelete := make(chan bool, 1)
	doneReadErrors := make(chan bool, 1)
	successCount := 0
	errorCount := 0

	deleteProgress := c.Ui.Progress("Deleting users")

	deleteProgress.Start()

	// send config
	deleteConfigReq := &proto.DeleteRequest{
		Data: &proto.DeleteRequest_Config{
			Config: providerConfigs,
		},
	}

	if err = deleteClient.Send(deleteConfigReq); err != nil {
		return errors.Wrap(err, "cannot sent config")
	}

	// send users
	go func() {
		for _, user := range cmd.UserIds {
			req := &proto.DeleteRequest{
				Data: &proto.DeleteRequest_UserId{
					UserId: user,
				},
			}
			c.Log.Trace().Msg(fmt.Sprintf("Deleting user: %s", req))
			if err = deleteClient.Send(req); err != nil {
				c.Log.Error().Msg(err.Error())
			}
			successCount++
		}
		doneDelete <- true

	}()

	// receive errors
	go func() {
		for {
			res, err := deleteClient.Recv()
			if err == io.EOF {
				doneReadErrors <- true
				return
			}
			errorCount++
			if err != nil {
				c.Log.Error().Msg(err.Error())
			}
			if respErr := res.GetError(); respErr != nil {
				c.Log.Error().Msg(respErr.Message)
			}
		}
	}()

	<-doneDelete
	if err = deleteClient.CloseSend(); err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	<-doneReadErrors

	deleteProgress.Stop()
	c.Ui.Normal().WithTable("Status", "N° of users").
		WithTableRow("Succeeded", fmt.Sprintf("%d", successCount)).
		WithTableRow("Errored", fmt.Sprintf("%d", errorCount)).Do()

	c.Ui.Success().Msg("Delete done")
	return nil
}
