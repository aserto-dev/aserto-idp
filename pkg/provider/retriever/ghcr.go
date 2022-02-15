package retriever

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/x"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/pkg/errors"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

var defaultRepoAddress = "ghcr.io/aserto-dev"

type GhcrRetriever struct {
	Store               *content.File
	RemoteStoreLocation string
	LocalStoreLocation  string
	extension           string
}

func NewGhcrRetriever() *GhcrRetriever {
	opSys := runtime.GOOS
	ext := ""
	if opSys == "windows" {
		ext = ".exe"
	}
	repoAdrress := os.Getenv("IDP_PLUGIN_REPO_ADDRESS")
	if repoAdrress == "" {
		repoAdrress = defaultRepoAddress
	}
	return &GhcrRetriever{
		extension:           ext,
		RemoteStoreLocation: fmt.Sprintf("%s/aserto-idp-plugins_%s_%s", repoAdrress, opSys, runtime.GOARCH),
	}
}

func (o *GhcrRetriever) Connect() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	o.LocalStoreLocation = filepath.Join(homeDir, ".aserto", "idpplugins")
	err = os.MkdirAll(o.LocalStoreLocation, 0777)
	if err != nil {
		return err
	}

	file := content.NewFile(o.LocalStoreLocation)

	o.Store = file

	return nil
}

func (o *GhcrRetriever) Disconnect() {
}

func (o *GhcrRetriever) List() ([]string, error) {
	repoName := o.RemoteStoreLocation
	repo, err := name.NewRepository(repoName)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid repo name [%s]", repoName)
	}

	tags, err := remote.List(repo)

	if err != nil {
		if tErr, ok := err.(*transport.Error); ok {
			switch tErr.StatusCode {
			case http.StatusUnauthorized:
				return nil, errors.Wrap(err, "authentication to docker registry failed")
			case http.StatusNotFound:
				return []string{}, nil
			}
		}

		return nil, errors.Wrap(err, "failed to list tags from registry")
	}

	return tags, nil
}

func (o *GhcrRetriever) Download(pluginName, version string) error {
	vers := strings.Split(version, ".")
	if vers[0] != IdpMajVersion() {
		return errors.New("incompatible version was provided for download; abort...") //nolint : revive : tbd
	}

	plgName := x.PluginPrefix + pluginName + o.extension
	destFilePath := filepath.Join(o.LocalStoreLocation, plgName)
	_, err := os.Stat(destFilePath)
	if err == nil {
		er := os.Remove(destFilePath)
		if er != nil {
			return errors.Wrap(err, "failed to remove old binary file")
		}
	}

	ref := fmt.Sprintf("%s:%s-%s", o.RemoteStoreLocation, pluginName, version)
	err = o.pull(ref)
	if err != nil {
		return err
	}

	err = os.Chmod(destFilePath, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to provide rights to output file [%s]", destFilePath)
	}

	return nil
}

func (o *GhcrRetriever) pull(ref string) error {
	resolver := docker.NewResolver(docker.ResolverOptions{
		Hosts: func(s string) ([]docker.RegistryHost, error) {
			client := &http.Client{}

			return []docker.RegistryHost{
				{
					Host:         s,
					Scheme:       "https",
					Capabilities: docker.HostCapabilityPull | docker.HostCapabilityResolve | docker.HostCapabilityPush,
					Client:       client,
					Path:         "/v2",
					Authorizer: docker.NewDockerAuthorizer(
						docker.WithAuthClient(client)),
				},
			}, nil
		},
	})

	allowedMediaTypes := []string{"application/vnd.unknown.layer.v1+txt", "application/vnd.unknown.config.v1+json"}
	opts := []oras.CopyOpt{
		oras.WithAllowedMediaTypes(allowedMediaTypes),
		oras.WithAdditionalCachedMediaTypes(allowedMediaTypes...),
	}
	_, err := oras.Copy(context.Background(), resolver, ref, o.Store, "", opts...)
	if err != nil {
		return errors.Wrapf(err, "download for '%s' failed", ref)
	}

	return nil
}
