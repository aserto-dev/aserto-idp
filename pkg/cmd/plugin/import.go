package plugin

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
)

type ImportCmd struct {
}

func (cmd *ImportCmd) Run(app *kong.Kong, context *kong.Context, c *cc.CC) error {
	configs, err := buildStructPb(context)
	if err != nil {
		return err
	}

	req := &proto.ImportRequest{
		Config: configs,
	}

	importClient, err := c.CommandIDP.Import(c.Context)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	return importClient.Send(req)
}
