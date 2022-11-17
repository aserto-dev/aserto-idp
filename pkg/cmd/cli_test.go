package cmd_test

import (
	"testing"

	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	logger "github.com/aserto-dev/aserto-logger"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	"github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var (
	checkForUpdates  = cmd.CheckForUpdates
	downloadProvider = cmd.DownloadProvider
	validatePlugin   = cmd.ValidatePlugin
)

func TestInvalidDownloadProvider(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	c.Log = &zerolog.Logger{}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	assert.NoError(err)

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)

	err = downloadProvider("plugin-name", c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "plugin 'plugin-name' does not exists")
}

func TestDownloadProvider(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("okta").AnyTimes()
	err = c.AddProvider(provider)
	assert.NoError(err)
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.1.0").Return(nil)

	err = downloadProvider("okta", c)
	assert.NoError(err)
}

func TestCheckForUpdatesNoUpdates(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	provider := mocks.NewMockProvider(ctrl)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &idpplugin.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.3.0"}, Description: "test"}

	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)
	provider.EXPECT().GetName().Return("auth0")
	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)

	updates, latestVer, err := checkForUpdates(provider, c)
	assert.NoError(err)
	assert.False(updates)
	assert.Empty(latestVer)
}

func TestCheckForUpdatesWithUpdates(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	provider := mocks.NewMockProvider(ctrl)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &idpplugin.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.1.0"}, Description: "test"}

	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)
	provider.EXPECT().GetName().Return("auth0")
	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)

	updates, latestVer, err := checkForUpdates(provider, c)
	assert.NoError(err)
	assert.True(updates)
	assert.Equal(latestVer, "0.2.0")
}
