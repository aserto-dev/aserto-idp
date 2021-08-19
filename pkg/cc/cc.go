package cc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aserto-dev/aserto-idp/pkg/auth0/api"
	"github.com/aserto-dev/aserto-idp/pkg/grpcc"
	"github.com/aserto-dev/aserto-idp/pkg/keyring"
	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/aserto-dev/aserto-idp/shared"
	"github.com/aserto-dev/go-lib/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// CC contains dependencies that are cross cutting and are needed in most
// of the providers that make up this application
type CC struct {
	Context     context.Context
	Log         *zerolog.Logger
	Plugin      shared.Provider
	services    *grpcc.Services
	overrides   map[string]string
	environment string
	Provider    string
	_token      *api.Token
}

func (ctx *CC) SetEnv(env string) error {
	log.Printf("set-context-env %s", env)
	if env == "" {
		return errors.Errorf("env is not set")
	}

	var err error
	ctx.services, err = grpcc.Environment(env)
	if err != nil {
		return err
	}

	ctx.environment = env

	return nil
}

func (ctx *CC) SetLogger(w io.Writer) {
	log.SetOutput(w)
}

func (ctx *CC) Override(key, value string) {
	log.Println("override-context-env", key, value)
	ctx.overrides[key] = value
}

func (ctx *CC) SetPlugin(plugin shared.Provider) {
	ctx.Plugin = plugin
}

func (ctx *CC) SetProvider(provider string) {
	ctx.Provider = provider
}

func (ctx *CC) AccessToken() string {
	return ctx.token().Access
}

func (ctx *CC) Token() *api.Token {
	return ctx.token()
}

func (ctx *CC) ExpiresAt() time.Time {
	return ctx.token().ExpiresAt
}

func (ctx *CC) TenantID() string {
	if tenantID, ok := ctx.overrides[x.TenantIDOverride]; ok {
		ctx.Log.Debug().Msg(fmt.Sprintf("!!! tenant override [%s]\n", tenantID))
		return tenantID
	}
	return ctx.token().TenantID
}

func (ctx *CC) AuthorizerService() string {
	if authorizer, ok := ctx.overrides[x.AuthorizerOverride]; ok {
		ctx.Log.Debug().Msg(fmt.Sprintf("!!! authorizer override [%s]\n", authorizer))
		return authorizer
	}
	return ctx.services.AuthorizerService
}

func (ctx *CC) token() *api.Token {
	if ctx._token == nil {
		kr, err := keyring.NewKeyRing(ctx.environment)
		if err != nil {
			log.Printf("token: instantiating keyring, %s", err.Error())
			return nil
		}

		ctx._token, err = kr.GetToken()
		if err != nil {
			return nil
		}
	}
	return ctx._token
}

func New() *CC {
	cfg := logger.Config{}
	cfg.LogLevel = "trace"
	cfg.LogLevelParsed = zerolog.TraceLevel

	log, _ := logger.NewLogger(os.Stdout, &cfg)

	ctx := CC{
		Context:   context.Background(),
		Log:       log,
		overrides: make(map[string]string),
	}
	return &ctx
}
