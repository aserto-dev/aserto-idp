package retriever

import (
	"regexp"
	"sort"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/version"
	versioning "github.com/hashicorp/go-version"
)

type PluginsInfo struct {
	Info map[string]map[string][]string
}

func NewPluginsInfo(source Retriever) *PluginsInfo {
	content, err := source.List()
	pluginInfo := &PluginsInfo{}
	if err == nil {
		pluginInfo.Info = make(map[string]map[string][]string)
		pluginInfo.extractInfo(content)
		pluginInfo.sortVersions()
	}
	return pluginInfo
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

	if pi.Info[pluginName] == nil {
		pi.Info[pluginName] = make(map[string][]string)
	}

	pi.Info[pluginName][majorVer] = append(pi.Info[pluginName][majorVer], version)
	pi.Info[pluginName][majorVer] = Unique(pi.Info[pluginName][majorVer])
}

func (pi *PluginsInfo) sortVersions() {
	for _, majVers := range pi.Info {
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

func (pi *PluginsInfo) LatestVersion(plugin string) string {
	cliMajor := IdpMajVersion()

	availableVersions := pi.Info[plugin][cliMajor]

	if len(availableVersions) == 0 {
		return ""
	}

	return availableVersions[0]
}

func IdpMajVersion() string {
	cliVer := version.GetVer()
	vers := strings.Split(cliVer, ".")

	return strings.ReplaceAll(vers[0], "v", "")
}
