module github.com/aserto-dev/aserto-idp

go 1.17

// replace github.com/aserto-dev/mage-loot => ../mage-loot

require (
	github.com/alecthomas/kong v0.2.17
	github.com/aserto-dev/clui v0.1.8
	github.com/aserto-dev/go-grpc v0.0.8-0.20211007104643-202ec2b6a966
	github.com/aserto-dev/go-utils v0.2.3
	github.com/aserto-dev/idp-plugin-sdk v0.0.1
	github.com/aserto-dev/mage-loot v0.4.8
	github.com/containerd/containerd v1.5.8
	github.com/google/go-containerregistry v0.7.0
	github.com/hashicorp/go-plugin v1.4.3
	github.com/magefile/mage v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/auth0.v5 v5.19.2
	oras.land/oras-go v0.4.0
)

require (
	github.com/docker/cli v20.10.11+incompatible // indirect
	github.com/docker/docker v20.10.11+incompatible // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
)
