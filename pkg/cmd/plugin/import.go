package plugin

import (
	"fmt"
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/pkg/errors"
)

type ImportCmd struct {
}

func (cmd *ImportCmd) Run(app *kong.Kong, context *kong.Context, c *cc.CC) error {

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

	err = validatePlugin(providerClient, c, providerConfigs, providerName)
	if err != nil {
		return err
	}

	exportClient, err := providerClient.Export(c.Context, req)
	if err != nil {
		return err
	}

	defaultProviderName := c.GetDefaultProvider().GetName()
	defaultProviderConfigs, err := getPbStructForNode(c.Config.Plugins[defaultProviderName], context.Path[0].Node())
	if err != nil {
		return err
	}

	defaultProviderClient, err := c.GetDefaultProvider().PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(defaultProviderClient, c, defaultProviderConfigs, defaultProviderName)
	if err != nil {
		return err
	}

	importClient, err := defaultProviderClient.Import(c.Context)
	if err != nil {
		return err
	}

	users := make(chan *api.User, 10)
	doneImport := make(chan bool, 1)
	doneExport := make(chan bool, 1)
	doneReadErrors := make(chan bool, 1)
	successCount := 0
	errorCount := 0

	importProgress := c.Ui.Progress("Importing users")

	importProgress.Start()

	// send config
	importConfigReq := &proto.ImportRequest{
		Data: &proto.ImportRequest_Config{
			Config: defaultProviderConfigs,
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
			if !includeExt {
				user.Attributes = &api.AttrSet{}
				user.Applications = make(map[string]*api.AttrSet)
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
				log.Fatalf("cannot receive %v", err)
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
	c.Ui.Normal().WithTable("Status", "N° of users").
		WithTableRow("Succeeded", fmt.Sprintf("%d", successCount)).
		WithTableRow("Errored", fmt.Sprintf("%d", errorCount)).Do()

	c.Ui.Success().Msg("Import done")
	return nil
}
