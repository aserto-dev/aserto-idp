package srv

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/aserto-idp/plugins/json/pkg/config"
	grpcerr "github.com/aserto-dev/aserto-idp/shared/errors"
	"github.com/aserto-dev/aserto-idp/shared/pb"
	"github.com/aserto-dev/aserto-idp/shared/version"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type JsonPluginServer struct{}

func (s JsonPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(config.GetVersion),
		Description: "Json Plugin",
		Configs:     config.GetPluginConfig(),
	}

	return &response, nil
}

func (s JsonPluginServer) Import(srv proto.Plugin_ImportServer) error {
	done := make(chan bool, 1)
	errDone := make(chan bool, 1)
	errc := make(chan error, 1)
	config := &config.JsonConfig{}
	count := int32(0)

	go grpcerr.SendImportErrors(srv, errc, errDone)

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
				if cfg := req.GetConfig(); cfg != nil {
					configBytes, err := protojson.Marshal(cfg)
					if err != nil {
						errc <- err
					}

					err = json.Unmarshal(configBytes, config)
					if err != nil {
						errc <- err
					}
				}
			}

			if user := req.GetUser(); user != nil {
				if u := user.GetUser(); u != nil {
					if count > 0 {
						_, _ = users.Write([]byte(",\n"))
					}
					b, err := options.Marshal(u)
					if err != nil {
						errc <- err
					}

					if _, err := users.Write(b); err != nil {
						errc <- err
					}
					count++
				}
			}
		}
	}()
	// Wait for EOF
	<-done

	_, err := users.Write([]byte("\n]\n"))
	if err != nil {
		errc <- err
	}
	f, err := os.Create(config.File)
	if err != nil {
		errc <- err
	}
	w := bufio.NewWriter(f)
	_, err = users.WriteTo(w)
	if err != nil {
		errc <- err
	}
	w.Flush()

	close(errc)
	<-errDone

	return nil
}

// func (s pluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return fmt.Errorf("not implemented")
// }

// Validate that the folder of the given file is writable
func (JsonPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}

	config, err := config.NewConfig(req.Config)
	if err != nil {
		return response, status.Error(codes.InvalidArgument, "failed to parse config")
	}

	dir := filepath.Dir(config.File)

	info, err := os.Stat(dir)
	if err != nil {
		return response, status.Error(codes.NotFound, err.Error())
	}

	if !info.IsDir() {
		return response, status.Errorf(codes.InvalidArgument, "%s is not a directory", dir)
	}

	if runtime.GOOS == "windows" {
		if info.Mode().Perm()&(1<<(uint(7))) == 0 {
			return response, status.Errorf(codes.PermissionDenied, "cannot access %s", dir)
		}
	} else {
		err = unix.Access(dir, unix.W_OK)
		if err != nil {
			return response, status.Errorf(codes.PermissionDenied, "cannot access %s: %s", dir, err.Error())
		}
	}

	return response, nil
}

func (s JsonPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {

	config, err := config.NewConfig(req.Config)
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
				User: &proto.User{
					Data: &proto.User_User{
						User: &u,
					},
				},
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
