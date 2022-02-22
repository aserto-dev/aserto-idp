package cmd

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRunListPluginsRemoteFalse(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("test-plugin").AnyTimes()

	err = c.AddProvider(provider)
	assert.NoError(err)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &idpplugin.InfoResponse{Build: &info.BuildInfo{Commit: "commit"}, Description: "test"}
	listCmd := &ListPluginsCmd{Remote: false}

	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)
	provider.EXPECT().GetName().Return("test-plugin").AnyTimes()

	err = listCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}

func TestRunListPluginsRemoteTrue(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	listCmd := &ListPluginsCmd{Remote: true}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	err = listCmd.Run(&kong.Context{}, c)

	assert.NoError(err)
}
