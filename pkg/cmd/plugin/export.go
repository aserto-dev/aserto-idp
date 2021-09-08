package plugin

import (
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/proto"

	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type ExportCmd struct {
	InclUserExt bool
	Source      string
	kong.Plugins
}

func (cmd *ExportCmd) Run(c *cc.CC) error {

	req := &proto.ExportRequest{
		Options: map[string]string{
			"source": cmd.Source,
		},
	}

	authorizerService := "authorizer.eng.aserto.com:8443"
	apiKey := "xxxxxxxxxxxxxxxxxxxxxxxx"
	tenant := "asdfagsg"

	export, err := c.CommandIDP.Export(c.Context, req)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	importSrv, err := c.DefaultIDP.Import(c.Context)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	for {
		resp, err := export.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		log.Printf("Resp received: %s", resp.Data)
		switch u := resp.Data.(type) {
		case *proto.ExportResponse_User:
			{
				req := &proto.ImportRequest{
					Options: map[string]string{
						"authorizer":  authorizerService,
						"api_key":     apiKey,
						"tenant":      tenant,
						"include_ext": "false",
					},
					Data: &proto.ImportRequest_User{
						User: u.User,
					},
				}
				if err = importSrv.Send(req); err != nil {
					log.Println(err)
					return err
				}
			}
		case *proto.ExportResponse_UserExt:
			{

			}
		}
	}
	return nil
}
