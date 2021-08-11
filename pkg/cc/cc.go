package cc

import (
	"context"
	"os"

	"github.com/aserto-dev/go-lib/logger"
	"github.com/rs/zerolog"
)

// CC contains dependencies that are cross cutting and are needed in most
// of the providers that make up this application
type CC struct {
	Context context.Context
	Log     *zerolog.Logger
}

func New() *CC {
	cfg := logger.Config{}
	cfg.LogLevel = "trace"
	cfg.LogLevelParsed = zerolog.TraceLevel

	log, _ := logger.NewLogger(os.Stdout, &cfg)

	ctx := CC{
		Context: context.Background(),
		Log:     log,
	}
	return &ctx
}
