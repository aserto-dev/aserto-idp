package errors

import (
	"log"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func SendImportErrors(srv proto.Plugin_ImportServer, errc <-chan error, done chan<- bool) {
	for {
		e, more := <-errc
		if !more {
			// channel closed
			done <- true
			return
		}
		err := srv.Send(
			&proto.ImportResponse{
				Error: &status.Status{
					Message: e.Error(),
				},
			},
		)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func SendExportErrors(srv proto.Plugin_ExportServer, errc <-chan error, done chan<- bool) {
	for {
		e, more := <-errc
		if !more {
			// channel closed
			done <- true
			return
		}
		err := srv.Send(
			&proto.ExportResponse{
				Data: &proto.ExportResponse_Error{
					Error: &status.Status{
						Message: e.Error(),
					},
				},
			},
		)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
