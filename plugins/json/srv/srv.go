package srv

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/json/config"
	"github.com/aserto-dev/aserto-idp/shared/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type JsonPluginServer struct{}

func (s JsonPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{}
	response.Build = "placeholder"
	response.System = ""
	response.Version = "placeholder"
	response.Description = "Json Plugin"
	response.Configs = config.GetPluginConfig()

	return &response, nil
}

func (s JsonPluginServer) Import(srv proto.Plugin_ImportServer) error {
	done := make(chan bool, 1)
	errc := make(chan error, 1)
	config := &config.JsonConfig{}
	count := int32(0)

	go func() {
		for e := range errc {
			log.Println(e.Error())
		}
	}()

	var users bytes.Buffer
	users.Write([]byte("[\n"))

	options := protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}

	go func() {
		for {
			req, err := srv.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				errc <- errors.Wrapf(err, "srv.Recv()")
			}

			if config.File == "" {
				configBytes, err := protojson.Marshal(req.Config)
				if err != nil {
					errc <- err
				}

				err = json.Unmarshal(configBytes, config)
				if err != nil {
					errc <- err
				}
			}

			switch u := req.Data.(type) {
			case *proto.ImportRequest_User:
				{
					if count > 0 {
						_, _ = users.Write([]byte(",\n"))
					}
					b, err := options.Marshal(u.User)
					if err != nil {
						errc <- err
					}

					if _, err := users.Write(b); err != nil {
						errc <- err
					}
					count++
				}
			case *proto.ImportRequest_UserExt:
				{

				}
			}
		}
	}()
	// Wait for EOF
	<-done

	res := &proto.ImportResponse{}
	res.SuccededCount = count
	res.FailCount = 0

	_, _ = users.Write([]byte("\n]\n"))
	f, _ := os.Create(config.File)
	w := bufio.NewWriter(f)
	users.WriteTo(w)
	w.Flush()
	err := srv.SendAndClose(res)
	if err != nil {
		errc <- err
	}

	// close error channel as the last action before returning
	close(errc)

	return nil
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// func (*pluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	return nil, fmt.Errorf("not implemented")
// }

func (s JsonPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {

	configBytes, err := protojson.Marshal(req.Config)
	if err != nil {
		return err
	}

	config := &config.JsonConfig{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return err
	}

	r, err := os.Open(config.File)
	if err != nil {
		log.Println(err)
		return err
	}

	dec := json.NewDecoder(r)

	if _, err = dec.Token(); err != nil {
		log.Println(err)
		return err
	}

	for dec.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(dec, &u); err != nil {
			log.Println(err)
			return err
		}
		res := &proto.ExportResponse{
			Data: &proto.ExportResponse_User{
				User: &u,
			},
		}
		if err = srv.Send(res); err != nil {
			log.Println(err)
			return err
		}
	}

	if _, err = dec.Token(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
