//go:build wireinject
// +build wireinject

package cc

import (
	"context"
	"io"

	"github.com/aserto-dev/aserto-idp/pkg/cc/config"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	"github.com/aserto-dev/aserto-idp/pkg/provider/retriever"
	"github.com/golang/mock/gomock"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/google/wire"
)

func BuildCC(logOutput logger.Writer, errorOutput logger.ErrWriter, output io.Writer, cfg *logger.Config) (*CC, error) {
	wire.Build(wire.Struct(new(CC), "Context", "Retriever", "Log", "UI", "Config", "pluginsInfo"),
		context.Background,
		logger.NewLogger,
		clui.NewUIWithOutput,
		wire.Bind(new(retriever.Retriever), new(*retriever.GhcrRetriever)),
		retriever.NewGhcrRetriever,
		config.NewEmptyConfig,
		retriever.NewPluginsInfo,
	)
	return &CC{}, nil
}

func BuildTestCC(ctrl *gomock.Controller, cfg *logger.Config) (*CC, error) {
	wire.Build(wire.Struct(new(CC), "Context", "Retriever", "Config", "pluginsInfo"),
		context.Background,
		wire.Bind(new(retriever.Retriever), new(*mocks.MockRetriever)),
		mocks.NewMockRetriever,
		config.NewEmptyConfig,
		retriever.NewPluginsInfo,
	)
	return &CC{}, nil
}
