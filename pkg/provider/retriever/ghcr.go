package retriever

import (
	"context"
	"fmt"
	"io"
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
	"github.com/rs/zerolog/log"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
)

type GhcrRetriever struct {
	Store               *content.OCIStore
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
	return &GhcrRetriever{
		extension:           ext,
		RemoteStoreLocation: fmt.Sprintf("ghcr.io/aserto-dev/aserto-idp-plugins_%s_%s", opSys, runtime.GOARCH),
	}
}

func (o *GhcrRetriever) Connect() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	o.LocalStoreLocation = filepath.Join(homeDir, ".aserto", "idpplugins", "ociStore")
	err = os.MkdirAll(o.LocalStoreLocation, 0777)
	if err != nil {
		return err
	}

	ociStore, err := content.NewOCIStore(o.LocalStoreLocation)
	if err != nil {
		return err
	}

	err = ociStore.LoadIndex()
	if err != nil {
		return err
	}

	o.Store = ociStore

	return nil
}

func (o *GhcrRetriever) Disconnect() {
	_ = os.RemoveAll(o.LocalStoreLocation)
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

func (o *GhcrRetriever) Download(pluginName string, version string) error {

	if version == "latest" {
		latestVersion := LatestVersion(pluginName, o)
		if latestVersion == "" {
			return fmt.Errorf("couldn't find latest version for %s", pluginName)
		}

		version = latestVersion
	}

	vers := strings.Split(version, ".")
	if vers[0] != IdpMajVersion() {
		return errors.New("incompatible version was provided for download; abort...")
	}

	ref := fmt.Sprintf("%s:%s-%s", o.RemoteStoreLocation, pluginName, version)
	err := o.pull(ref)
	if err != nil {
		return err
	}

	dest := strings.ReplaceAll(o.LocalStoreLocation, "ociStore", "")
	destFilePath := filepath.Join(dest, x.PluginPrefix+pluginName+o.extension)
	err = o.save(ref, destFilePath)
	if err != nil {
		return err
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
	opts := []oras.PullOpt{
		oras.WithAllowedMediaTypes(allowedMediaTypes),
		oras.WithCachedMediaTypes(allowedMediaTypes...),
		oras.WithContentProvideIngester(o.Store),
	}
	_, descriptors, err := oras.Pull(context.Background(), resolver, ref, o.Store, opts...)
	if err != nil {
		return errors.Wrapf(err, "download for %s failed", ref)
	}

	if len(descriptors) != 1 {
		return errors.Errorf("unexpected layer count of [%d] from the registry; expected 1", len(descriptors))
	}

	o.Store.AddReference(ref, descriptors[0])
	err = o.Store.SaveIndex()
	if err != nil {
		return err
	}

	return nil
}

func (o *GhcrRetriever) save(ref, outputFile string) error {
	err := o.Store.LoadIndex()
	if err != nil {
		return err
	}

	refs := o.Store.ListReferences()

	refDescriptor, ok := refs[ref]
	if !ok {
		return errors.Errorf("provider [%s] not found in the local store", ref)
	}
	reader, err := o.Store.ReaderAt(context.Background(), refDescriptor)
	if err != nil {
		return errors.Wrap(err, "failed to open store reader")
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	out, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrapf(err, "failed to create output file [%s]", outputFile)
	}
	err = os.Chmod(outputFile, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to provide rights to output file [%s]", outputFile)
	}

	defer func() {
		err := out.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	chunkSize := 64
	buf := make([]byte, chunkSize)
	for i := 0; i < int(reader.Size()); {
		if chunkSize < int(reader.Size())-i {
			chunkSize = int(reader.Size()) - i
			buf = make([]byte, chunkSize)
		}

		n, err := reader.ReadAt(buf, int64(i))
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "failed to read OCI idp binary")
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			return errors.Wrap(err, "failed to write idp binary to file")
		}

		if err == io.EOF {
			break
		}

		i += chunkSize
	}

	return nil
}
