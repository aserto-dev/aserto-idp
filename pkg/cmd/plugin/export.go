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

type ExportCmd struct {
}

func (cmd *ExportCmd) Run(app *kong.Kong, context *kong.Context, c *cc.CC) error {

	users := make(chan *api.User, 10)
	doneImport := make(chan bool, 1)
	doneExport := make(chan bool, 1)
	doneReadErrors := make(chan bool, 1)
	successCount := 0
	errorCount := 0

	defaultProviderName := c.GetDefaultProvider().GetName()
	defaultProviderConfigs, err := getPbStructForNode(c.Config.Plugins[defaultProviderName], context.Path[0].Node())
	if err != nil {
		return err
	}

	//TODO: Better handle this
	includeExt := false

	exReq := &proto.ExportRequest{
		Config: defaultProviderConfigs,
	}

	defaultProviderClient, err := c.GetDefaultProvider().PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(defaultProviderClient, c, defaultProviderConfigs, defaultProviderName)
	if err != nil {
		return err
	}

	exportClient, err := defaultProviderClient.Export(c.Context, exReq)
	if err != nil {
		return err
	}

	providerName := context.Selected().Parent.Name

	configs, err := getPbStructForNode(c.Config.Plugins[providerName], context.Selected().Parent)
	if err != nil {
		return err
	}

	providerClient, err := c.GetProvider(providerName).PluginClient()
	if err != nil {
		return err
	}

	err = validatePlugin(providerClient, c, configs, providerName)
	if err != nil {
		return err
	}

	importClient, err := providerClient.Import(c.Context)
	if err != nil {
		return err
	}

	// send config
	req := &proto.ImportRequest{
		Data: &proto.ImportRequest_Config{
			Config: configs,
		},
	}

	if err = importClient.Send(req); err != nil {
		return errors.Wrap(err, "cannot sent config")
	}

	exportProgress := c.Ui.Progress("Exporting users")

	exportProgress.Start()
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
				continue
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
	exportProgress.Stop()

	c.Ui.Normal().WithTable("Status", "N° of users").
		WithTableRow("Succeeded", fmt.Sprintf("%d", successCount)).
		WithTableRow("Errored", fmt.Sprintf("%d", errorCount)).Do()

	c.Ui.Success().Msg("Export done")
	return nil
}
