package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRunDeleteWithoutFromPlugin(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	deleteCmd := &DeleteCmd{}

	err = deleteCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "rpc error: code = InvalidArgument desc = no '--from' idp was provided")
}

func TestRunDeleteWithoutPluginsOnTheSystem(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	deleteCmd := &DeleteCmd{From: "okta", NoUpdateCheck: true}

	err = deleteCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Contains(err.Error(), "unavailable \"--from\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
}

func TestRunDeleteWithoutUserIds(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	kongFlag := &kong.Flag{Group: nil}
	kongPath := &kong.Path{App: &kong.Application{Node: &kong.Node{Flags: []*kong.Flag{kongFlag}}}}
	kongContext := &kong.Context{Path: []*kong.Path{kongPath}}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("okta").AnyTimes()
	deleteClient := mocks.NewMockPlugin_DeleteClient(ctrl)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.1.0"}, Description: "test"}

	err = c.AddProvider(provider)
	assert.NoError(err)
	deleteCmd := &DeleteCmd{From: "okta"}

	provider.EXPECT().PluginClient().Return(pluginClient, nil).Times(2)
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil).AnyTimes()
	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	pluginClient.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil, nil)
	pluginClient.EXPECT().Delete(gomock.Any()).Return(deleteClient, nil)
	deleteClient.EXPECT().Send(gomock.Any()).Return(nil)
	deleteClient.EXPECT().Recv().Return(nil, io.EOF)
	deleteClient.EXPECT().CloseSend().Return(nil)

	err = deleteCmd.Run(kongContext, c)
	assert.NoError(err)
}

func TestRunDeleteWithUserIds(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	kongFlag := &kong.Flag{Group: nil}
	kongPath := &kong.Path{App: &kong.Application{Node: &kong.Node{Flags: []*kong.Flag{kongFlag}}}}
	kongContext := &kong.Context{Path: []*kong.Path{kongPath}}
	pluginsVersions := []string{"okta-0.1.0", "auth0-0.2.0"}
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("okta").AnyTimes()
	deleteClient := mocks.NewMockPlugin_DeleteClient(ctrl)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	deleteReq := &proto.DeleteRequest{
		Data: &proto.DeleteRequest_UserId{
			UserId: "someUserID",
		},
	}
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.1.0"}, Description: "test"}

	err = c.AddProvider(provider)
	assert.NoError(err)
	deleteCmd := &DeleteCmd{From: "okta", UserIds: []string{"someUserID"}}

	provider.EXPECT().PluginClient().Return(pluginClient, nil).AnyTimes()
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil).AnyTimes()
	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil).AnyTimes()
	pluginClient.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	pluginClient.EXPECT().Delete(gomock.Any()).Return(deleteClient, nil)
	deleteClient.EXPECT().Send(deleteReq).Return(nil)
	deleteClient.EXPECT().Send(gomock.Any()).Return(nil)
	deleteClient.EXPECT().Recv().Return(nil, io.EOF)
	deleteClient.EXPECT().CloseSend().Return(nil)

	err = deleteCmd.Run(kongContext, c)
	assert.NoError(err)
}
