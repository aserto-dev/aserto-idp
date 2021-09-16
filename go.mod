module github.com/aserto-dev/aserto-idp

go 1.16

// replace github.com/aserto-dev/go-utils => ../go-utils

require (
	github.com/alecthomas/kong v0.2.17
	github.com/aserto-dev/go-grpc v0.0.3
	github.com/aserto-dev/go-grpc-authz v0.0.2
	github.com/aserto-dev/go-utils v0.0.4
	github.com/aserto-dev/mage-loot v0.4.7
	github.com/hashicorp/go-plugin v1.4.2
	github.com/jhump/protoreflect v1.9.0 // indirect
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210811021853-ddbe55d93216
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
