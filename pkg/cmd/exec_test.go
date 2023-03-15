package cmd_test

import (
	"io"
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
	"github.com/aserto-dev/aserto-idp/pkg/cmd"
	"github.com/aserto-dev/aserto-idp/pkg/mocks"
	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/go-grpc/aserto/common/info/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/logger"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRunExecWithoutFromPlugin(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	execCmd := &cmd.ExecCmd{To: "json"}

	err = execCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "rpc error: code = InvalidArgument desc = no '--from' idp was provided")
}

func TestRunExecWithoutToPlugin(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	execCmd := &cmd.ExecCmd{From: "json"}

	err = execCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "rpc error: code = InvalidArgument desc = no '--to' idp was provided")
}

func TestRunExecWithoutPluginsOnTheSystem(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	execCmd := &cmd.ExecCmd{From: "json", To: "okta", NoUpdateCheck: true}

	err = execCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Contains(err.Error(), "unavailable \"--from\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
}

func TestRunExecWithoutToPluginOnTheSystem(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	provider := mocks.NewMockProvider(ctrl)
	provider.EXPECT().GetName().Return("json").AnyTimes()

	err = c.AddProvider(provider)
	assert.NoError(err)
	execCmd := &cmd.ExecCmd{From: "json", To: "okta", NoUpdateCheck: true}

	err = execCmd.Run(&kong.Context{}, c)
	assert.NotNil(err)
	assert.Contains(err.Error(), "unavailable \"--to\" idp was provided, use exec without --no-update-check to download it or use get-plugin command")
}

func TestRunExec(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	c, err := cc.BuildTestCC(ctrl, &logger.Config{})
	assert.NoError(err)
	c.Log = &zerolog.Logger{}
	c.UI = clui.NewUIWithOutput(os.Stdout)
	fromProvider := mocks.NewMockProvider(ctrl)
	fromProvider.EXPECT().GetName().Return("okta").AnyTimes()
	toProvider := mocks.NewMockProvider(ctrl)
	toProvider.EXPECT().GetName().Return("json").AnyTimes()
	kongFlag := &kong.Flag{Group: nil}
	kongPath := &kong.Path{App: &kong.Application{Node: &kong.Node{Flags: []*kong.Flag{kongFlag}}}}
	kongContext := &kong.Context{Path: []*kong.Path{kongPath}}
	exportClient := mocks.NewMockPlugin_ExportClient(ctrl)
	importClient := mocks.NewMockPlugin_ImportClient(ctrl)
	pluginClient := mocks.NewMockPluginClient(ctrl)
	providerInfo := &proto.InfoResponse{Build: &info.BuildInfo{Commit: "commit", Version: "v0.1.0"}, Description: "test"}
	pluginsVersions := []string{"okta-0.1.0", "json-0.1.0"}

	err = c.AddProvider(fromProvider)
	assert.NoError(err)
	err = c.AddProvider(toProvider)
	assert.NoError(err)
	execCmd := &cmd.ExecCmd{From: "okta", To: "json"}

	fromProvider.EXPECT().PluginClient().Return(pluginClient, nil).AnyTimes()
	pluginClient.EXPECT().Info(gomock.Any(), gomock.Any()).Return(providerInfo, nil).AnyTimes()
	c.Retriever.(*mocks.MockRetriever).EXPECT().List().Return(pluginsVersions, nil)
	toProvider.EXPECT().PluginClient().Return(pluginClient, nil).AnyTimes()
	pluginClient.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
	pluginClient.EXPECT().Export(gomock.Any(), gomock.Any()).Return(exportClient, nil)
	pluginClient.EXPECT().Import(gomock.Any(), gomock.Any()).Return(importClient, nil)
	importClient.EXPECT().Send(gomock.Any()).Return(nil)
	exportClient.EXPECT().Recv().Return(nil, io.EOF)
	importClient.EXPECT().CloseSend().Return(nil)
	importClient.EXPECT().Recv().Return(nil, io.EOF)

	err = execCmd.Run(kongContext, c)
	assert.NoError(err)
}
