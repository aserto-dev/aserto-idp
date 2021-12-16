package retriever

import (
	"regexp"
	"sort"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/version"
	versioning "github.com/hashicorp/go-version"
)

type PluginsInfo struct {
	info         map[string]map[string][]string
	RemoteSource Retriever
}

func NewPluginsInfo(source Retriever) *PluginsInfo {
	return &PluginsInfo{RemoteSource: source}
}

func (pi *PluginsInfo) extractInfo(versions []string) {
	nameRegex, err := regexp.Compile("([a-z0-9]+)-")
	if err != nil {
		return
	}
	versionRegex, err := regexp.Compile(`([0-9]+\.)([0-9]+\.)([0-9]+)`)
	if err != nil {
		return
	}

	for _, version := range versions {
		mat := nameRegex.FindStringSubmatch(version)
		if len(mat) == 0 {
			continue
		}

		name := mat[0][0 : len(mat[0])-1]

		mat = versionRegex.FindStringSubmatch(version)

		if len(mat) == 0 {
			continue
		}

		version := mat[0]

		reMajor := strings.Split(version, ".")
		major := reMajor[0]

		pi.populateVersions(name, major, version)
	}
}

func (pi *PluginsInfo) populateVersions(pluginName, majorVer, version string) {

	if pi.info[pluginName] == nil {
		pi.info[pluginName] = make(map[string][]string)
	}

	pi.info[pluginName][majorVer] = append(pi.info[pluginName][majorVer], version)
	pi.info[pluginName][majorVer] = Unique(pi.info[pluginName][majorVer])
}

func (pi *PluginsInfo) sortVersions() {
	for _, majVers := range pi.info {
		for _, vers := range majVers {
			length := len(vers)
			var objectVers []*versioning.Version
			for _, ver := range vers {
				v, err := versioning.NewVersion(ver)
				if err != nil {
					break
				}
				objectVers = append(objectVers, v)
			}
			if length != len(objectVers) {
				continue
			}

			sort.Slice(objectVers, func(i, j int) bool {
				return objectVers[i].GreaterThan(objectVers[j])
			})

			for index, ver := range objectVers {
				vers[index] = ver.String()
			}
		}

	}
}

func (pi *PluginsInfo) initializeInfo() error {
	content, err := pi.RemoteSource.List()
	if err != nil {
		return err
	}

	pi.info = make(map[string]map[string][]string)
	pi.extractInfo(content)
	pi.sortVersions()
	return nil
}

func (pi *PluginsInfo) GetInfo() (map[string]map[string][]string, error) {
	if pi.info == nil {
		err := pi.initializeInfo()
		return pi.info, err
	}

	return pi.info, nil
}

func (pi *PluginsInfo) LatestVersion(plugin string) (string, error) {
	cliMajor := IdpMajVersion()

	if pi.info == nil {
		err := pi.initializeInfo()
		if err != nil {
			return "", err
		}
	}

	availableVersions := pi.info[plugin][cliMajor]

	if len(availableVersions) == 0 {
		return "", nil
	}

	return availableVersions[0], nil
}

func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func IdpMajVersion() string {
	cliVer := version.GetVer()
	vers := strings.Split(cliVer, ".")

	return strings.ReplaceAll(vers[0], "v", "")
}
