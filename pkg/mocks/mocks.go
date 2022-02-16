package mocks

//go:generate mockgen -destination=mock_retriever.go -package=mocks github.com/aserto-dev/aserto-idp/pkg/provider/retriever Retriever
//go:generate mockgen -destination=mock_provider.go -package=mocks github.com/aserto-dev/aserto-idp/pkg/provider Provider
//go:generate mockgen -destination=mock_plugin_client.go -package=mocks github.com/aserto-dev/idp-plugin-sdk/grpcplugin PluginClient
//go:generate mockgen -destination=mock_servers.go -package=mocks github.com/aserto-dev/go-grpc/aserto/idpplugin/v1 Plugin_DeleteClient,Plugin_ExportClient,Plugin_ImportClient
