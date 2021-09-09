package plugin

import (
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

type ExportCmd struct {
}

func (cmd *ExportCmd) Run(app *kong.Kong, context *kong.Context, c *cc.CC) error {

	configs, err := buildStructPb(context)
	if err != nil {
		return err
	}

	req := &proto.ExportRequest{
		Config: configs,
	}

	authorizerService := "authorizer.eng.aserto.com:8443"
	apiKey := "XXXX"
	tenant := "XXXX"
	includeExt := false

	exportClient, err := c.CommandIDP.Export(c.Context, req)
	if err != nil {
		c.Log.Debug().Msg(err.Error())
	}

	importClient, err := c.DefaultIDP.Import(c.Context)
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

			configs := map[string]interface{}{
				"authorizer":  authorizerService,
				"api_key":     apiKey,
				"tenant":      tenant,
				"include_ext": false,
			}

			configStruct, err := structpb.NewStruct(configs)
			if err != nil {
				errc <- errors.Wrapf(err, "could not create config struct")
			}

			req := &proto.ImportRequest{
				Config: configStruct,
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
