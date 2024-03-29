package cmd

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExecCmd struct {
	From          string `short:"f" help:"The idp name you want to import from"`
	To            string `short:"t" help:"The idp name you want to import to"`
	NoUpdateCheck bool   `short:"n" help:"Don't check for plugins updates"`
}

func (cmd *ExecCmd) Run(context *kong.Context, c *cc.CC) error { //nolint : funlen // tbd

	if cmd.From == "" {
		return status.Error(codes.InvalidArgument, "no '--from' idp was provided")
	}

	if cmd.To == "" {
		return status.Error(codes.InvalidArgument, "no '--to' idp was provided")
	}

	if c.ProviderExists(cmd.From) && c.ProviderExists(cmd.To) && !cmd.NoUpdateCheck {
		sourceUpdates, latestFrom, err := checkForUpdates(c.GetProvider(cmd.From), c)
		if err != nil {
			c.UI.Exclamation().WithErr(err).Msgf("failed to check for updates on '%s' plugin", cmd.From)
		}
		destinationUpdates, latestTo, err := checkForUpdates(c.GetProvider(cmd.To), c)
		if err != nil {
			c.UI.Exclamation().WithErr(err).Msgf("failed to check for updates on '%s' plugin", cmd.To)
		}

		if sourceUpdates {
			c.UI.Exclamation().Msgf("a new version '%s' of '%s' plugin is available", latestFrom, cmd.From)
		}
		if destinationUpdates {
			c.UI.Exclamation().Msgf("a new version '%s' of '%s' plugin is available", latestTo, cmd.To)
		}
	}

	if !c.ProviderExists(cmd.From) {
		if cmd.NoUpdateCheck {
			return status.Error(codes.InvalidArgument, "unavailable \"--from\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
		}
		err := downloadProvider(cmd.From, c)
		if err != nil {
			return err
		}
	}

	if !c.ProviderExists(cmd.To) {
		if cmd.NoUpdateCheck {
			return status.Error(codes.InvalidArgument, "unavailable \"--to\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
		}
		err := downloadProvider(cmd.To, c)
		if err != nil {
			return err
		}
	}

	sourceProviderName := cmd.From
	sourceProviderConfigs, err := getPbStructForNode(sourceProviderName, c.Config.Plugins[sourceProviderName], context.Path[0].Node())
	if err != nil {
		return err
	}
	req := &proto.ExportRequest{
		Config: sourceProviderConfigs,
	}

	sourceProviderClient, err := c.GetProvider(sourceProviderName).PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(sourceProviderClient, c, sourceProviderConfigs, sourceProviderName, proto.OperationType_OPERATION_TYPE_EXPORT)
	if err != nil {
		return err
	}

	exportClient, err := sourceProviderClient.Export(c.Context, req)
	if err != nil {
		return err
	}

	destinationProviderName := cmd.To
	destinationProviderConfigs, err := getPbStructForNode(destinationProviderName, c.Config.Plugins[destinationProviderName], context.Path[0].Node())
	if err != nil {
		return err
	}

	destinationProviderClient, err := c.GetProvider(destinationProviderName).PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(destinationProviderClient, c, destinationProviderConfigs, destinationProviderName, proto.OperationType_OPERATION_TYPE_IMPORT)
	if err != nil {
		return err
	}

	importClient, err := destinationProviderClient.Import(c.Context)
	if err != nil {
		return err
	}

	users := make(chan *api.User, 10)
	doneImport := make(chan bool, 1)
	doneExport := make(chan bool, 1)
	doneReadErrors := make(chan bool, 1)
	successCount := 0
	errorCount := 0

	importProgress := c.UI.Progress("Importing users")

	importProgress.Start()

	// send config
	importConfigReq := &proto.ImportRequest{
		Data: &proto.ImportRequest_Config{
			Config: destinationProviderConfigs,
		},
	}

	if err = importClient.Send(importConfigReq); err != nil {
		return errors.Wrap(err, "cannot sent config")
	}

	// send users
	go func() {
		for {
			user, more := <-users
			if !more {
				doneImport <- true
				return
			}

			req := &proto.ImportRequest{
				Data: &proto.ImportRequest_User{
					User: user,
				},
			}
			c.Log.Trace().Msg(fmt.Sprintf("Sending user: %s", req))
			if err = importClient.Send(req); err != nil {
				c.Log.Error().Msg(err.Error())
			}
		}
	}()

	// receive errors
	go func() {
		for {
			res, err := importClient.Recv()
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

	// receive users
	go func() {
		for {
			resp, err := exportClient.Recv()
			if err == io.EOF {
				doneExport <- true
				return
			}
			if err != nil {
				c.Log.Error().Msgf("cannot receive %v", err)
				doneExport <- true
				return
			}
			c.Log.Trace().Msg(fmt.Sprintf("Resp received: %s", resp.Data))

			if respErr := resp.GetError(); respErr != nil {
				c.Log.Error().Msg(respErr.Message)
				continue
			}

			if user := resp.GetUser(); user != nil {
				users <- user
			}
			successCount++
		}
	}()

	<-doneExport
	close(users)
	<-doneImport
	if err = importClient.CloseSend(); err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	<-doneReadErrors

	importProgress.Stop()
	c.UI.Normal().WithTable("Status", "N° of users").
		WithTableRow("Succeeded", fmt.Sprintf("%d", successCount)).
		WithTableRow("Errored", fmt.Sprintf("%d", errorCount)).Do()

	c.UI.Success().Msg("Import done")
	return nil
}
