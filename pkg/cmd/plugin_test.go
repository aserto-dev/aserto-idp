package cmd_test

import (
	"os"
	"testing"

	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"

	logger "github.com/aserto-dev/aserto-logger"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewPlugin(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	provider := mocks.NewMockProvider(ctrl)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	c.Log = &zerolog.Logger{}
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit"}, Description: "test"}
	assert.NoError(err)

	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)
	provider.EXPECT().GetName().Return("test-name")

	plugin, err := cmd.NewPlugin(provider, c)
	assert.NoError(err)
	assert.NotNil(plugin)
	assert.Equal(plugin.Description, "test")
	assert.Equal(plugin.Name, "test-name")
}

func TestValidatePlugin(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	assert.NoError(err)

	pluginClient.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil, nil)

	err = validatePlugin(pluginClient, c, &structpb.Struct{}, "test-plugin", proto.OperationType_OPERATION_TYPE_EXPORT)
	assert.NoError(err)
}
