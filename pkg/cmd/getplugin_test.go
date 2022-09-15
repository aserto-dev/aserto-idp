package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	logger "github.com/aserto-dev/aserto-logger"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRunGetPluginInvalidName(t *testing.T) {
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
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "test-plugin"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "couldn't find latest version for test-plugin")
}

func TestRunGetPluginWithoutVersion(t *testing.T) {
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
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.1.0").Return(nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}

func TestRunGetPluginWithLatest(t *testing.T) {
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
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta:latest"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.1.0").Return(nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}

func TestRunGetPluginWithValidVersion(t *testing.T) {
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
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta:0.1.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.1.0").Return(nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}

func TestRunGetPluginWithInvalidVersion(t *testing.T) {
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
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta:0.3.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.3.0").Return(errors.New("Invalid version"))

	err = getCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "Invalid version")
}

func TestRunGetPluginAlreadyAtLatest(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("okta").AnyTimes()

	err = c.AddProvider(provider)
	assert.NoError(err)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "0.1.0"}, Description: "test"}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta:0.1.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}

func TestRunGetPluginNotAtLatest(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("okta").AnyTimes()

	err = c.AddProvider(provider)
	assert.NoError(err)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.1.0"}, Description: "test"}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	getCmd := &GetPluginCmd{Plugin: "okta:0.1.0"}

	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	provider.EXPECT().PluginClient().Return(pluginClient, nil)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil)
	provider.EXPECT().Kill()
	c.Retriever.(*mocks.MockRetriever).EXPECT().Download("okta", "0.1.0").Return(nil)

	err = getCmd.Run(&kong.Context{}, c)
	assert.NoError(err)
}
