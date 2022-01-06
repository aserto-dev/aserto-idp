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
	From          string   `short:"f" help:"The idp provider name you want to delete from"`
	UserIds       []string `arg:"" name:"user_id" help:"Users to remove." type:"string"`
	NoUpdateCheck bool     `short:"n" help:"Don't check for plugins updates"`
	kong.Plugins
}

func (cmd *DeleteCmd) Run(context *kong.Context, c *cc.CC) error {

	if cmd.From == "" {
		return status.Error(codes.InvalidArgument, "no '--from' idp was provided")
	}

	if !cmd.NoUpdateCheck && c.ProviderExists(cmd.From) {
		sourceUpdates, latest, err := checkForUpdates(c.GetProvider(cmd.From), c)
		if err != nil {
			c.Ui.Exclamation().WithErr(err).Msgf("Failed to check for updates for plugin '%s'", cmd.From)
		}

		if sourceUpdates {
			c.Ui.Exclamation().Msgf("A new version '%s' of the plugin '%s' is available", latest, cmd.From)
		}
	}

	if !c.ProviderExists(cmd.From) {
		if cmd.NoUpdateCheck {
			return status.Error(codes.InvalidArgument, "unavailable \"--from\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
		} else {
			err := downloadProvider(cmd.From, c)
			if err != nil {
				return err
			}
		}
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

	err = validatePlugin(providerClient, c, providerConfigs, providerName, proto.OperationType_OPERATION_TYPE_DELETE)
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
