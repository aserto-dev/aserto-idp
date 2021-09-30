module github.com/aserto-dev/aserto-idp

go 1.16

// replace github.com/aserto-dev/go-utils => ../go-utils

require (
	github.com/alecthomas/kong v0.2.17
	github.com/aserto-dev/clui v0.1.2
	github.com/aserto-dev/go-grpc v0.0.6
	github.com/aserto-dev/go-grpc-authz v0.0.2
	github.com/aserto-dev/go-utils v0.0.4
	github.com/aserto-dev/idp-plugin-sdk v0.0.1
	github.com/aserto-dev/mage-loot v0.4.8
	github.com/hashicorp/go-plugin v1.4.3
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/auth0.v5 v5.19.2
)
