package retriever

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/aserto-dev/aserto-idp/pkg/version"
	versioning "github.com/hashicorp/go-version"
)

var once sync.Once

type PluginsInfoSingleton map[string]map[string][]string

var (
	instance PluginsInfoSingleton
)

func PluginVersions(source Retriever) PluginsInfoSingleton {

	once.Do(func() {
		content, err := source.List()
		if err == nil {
			instance = make(PluginsInfoSingleton)
			ExtractInfo(content, instance)
			sortVersions(instance)
		}
	})

	return instance
}

func ExtractInfo(versions []string, singl PluginsInfoSingleton) {
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

		populateVersions(name, major, version, singl)
	}
}

func populateVersions(pluginName, majorVer, version string, singl PluginsInfoSingleton) {

	if singl[pluginName] == nil {
		singl[pluginName] = make(map[string][]string)
	}

	singl[pluginName][majorVer] = append(singl[pluginName][majorVer], version)
	singl[pluginName][majorVer] = Unique(singl[pluginName][majorVer])
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

func sortVersions(versions PluginsInfoSingleton) {
	for _, majVers := range versions {
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

func IdpMajVersion() string {
	cliVer := version.GetVer()
	vers := strings.Split(cliVer, ".")

	return strings.ReplaceAll(vers[0], "v", "")
}

func LatestVersion(plugin string, source Retriever) string {
	cliMajor := IdpMajVersion()

	pluginInfo := PluginVersions(source)

	availableVersions := pluginInfo[plugin][cliMajor]

	if len(availableVersions) == 0 {
		return ""
	}

	return availableVersions[0]
}
