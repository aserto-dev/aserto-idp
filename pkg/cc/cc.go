package cc

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aserto-dev/aserto-idp/shared/grpcplugin"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/rs/zerolog"
)

// CC contains dependencies that are cross cutting and are needed in most
// of the providers that make up this application
type CC struct {
	Context          context.Context
	Log              *zerolog.Logger
	DefaultIDPClient grpcplugin.PluginClient
	IDPClients       map[string]grpcplugin.PluginClient
	// services    *grpcc.Services
	// overrides   map[string]string
	// environment string
	// Provider    string
	// APIKey      string
}

// func (ctx *CC) SetEnv(env string) error {
// 	log.Printf("set-context-env %s", env)
// 	if env == "" {
// 		return errors.Errorf("env is not set")
// 	}

// 	var err error
// 	ctx.services, err = grpcc.Environment(env)
// 	if err != nil {
// 		return err
// 	}

// 	ctx.environment = env

// 	return nil
// }

func (ctx *CC) SetLogger(w io.Writer) {
	log.SetOutput(w)
}

// func (ctx *CC) Override(key, value string) {
// 	log.Println("override-context-env", key, value)
// 	ctx.overrides[key] = value
// }

// func (ctx *CC) SetAPIKey(key string) {
// 	ctx.APIKey = key
// }

// func (ctx *CC) AuthorizerService() string {
// 	if authorizer, ok := ctx.overrides[x.AuthorizerOverride]; ok {
// 		ctx.Log.Debug().Msg(fmt.Sprintf("!!! authorizer override [%s]\n", authorizer))
// 		return authorizer
// 	}
// 	return ctx.services.AuthorizerService
// }

func New() *CC {
	cfg := logger.Config{}
	cfg.LogLevel = "trace"
	cfg.LogLevelParsed = zerolog.TraceLevel

	log, _ := logger.NewLogger(os.Stdout, &cfg)

	ctx := CC{
		Context:    context.Background(),
		Log:        log,
		IDPClients: make(map[string]grpcplugin.PluginClient),
		// overrides: make(map[string]string),
	}
	return &ctx
}
