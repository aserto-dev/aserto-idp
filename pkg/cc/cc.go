package cc

import (
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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
}

func (ctx *CC) SetLogger(w io.Writer) {
	log.SetOutput(w)
}

func New() *CC {
	cfg := logger.Config{}

	cfg.LogLevelParsed = getLogLevel()

	log, _ := logger.NewLogger(os.Stdout, &cfg)

	ctx := CC{
		Context:    context.Background(),
		Log:        log,
		IDPClients: make(map[string]grpcplugin.PluginClient),
	}
	return &ctx
}

func getLogLevel() zerolog.Level {
	logLevel := zerolog.ErrorLevel

	for _, arg := range os.Args {
		if strings.HasPrefix(strings.ToLower(arg), "--verbosity=") {
			intLevel, err := strconv.Atoi(strings.Split(arg, "=")[1])
			if err != nil {
				break
			}
			switch intLevel {
			case 1:
				logLevel = zerolog.InfoLevel
			case 2:
				logLevel = zerolog.DebugLevel
			case 3:
				logLevel = zerolog.TraceLevel
			}

		}
		switch arg {
		case "-v":
			logLevel = zerolog.InfoLevel
		case "-vv":
			logLevel = zerolog.DebugLevel
		case "-vvv":
			logLevel = zerolog.TraceLevel
		}

	}

	return logLevel
}
