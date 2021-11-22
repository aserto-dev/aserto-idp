package retriever

import (
	"regexp"
	"sync"

	"github.com/aserto-dev/aserto-idp/pkg/version"
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
		}
	})

	return instance
}

func ExtractInfo(versions []string, singl PluginsInfoSingleton) {

	for _, version := range versions {
		re := regexp.MustCompile("([a-z0-9]+)-")
		mat := re.FindStringSubmatch(version)
		if len(mat) == 0 {
			continue
		}

		name := mat[0][0 : len(mat[0])-1]

		reVersion := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
		mat = reVersion.FindStringSubmatch(version)

		if len(mat) == 0 {
			continue
		}

		version := mat[0]

		reMajor := regexp.MustCompile(".[0-9]+.")
		mat = reMajor.FindStringSubmatch(version)

		if len(mat) == 0 {
			continue
		}

		major := mat[0][1 : len(mat[0])-1]

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

func IdpMajVersion() string {
	cliVer := version.GetVer()
	reMajor := regexp.MustCompile(".[0-9]+.")
	match := reMajor.FindStringSubmatch(cliVer)

	return match[0][1 : len(match[0])-1]
}

func LatestVersion(plugin string, source Retriever) string {
	cliVer := version.GetVer()
	reMajor := regexp.MustCompile(".[0-9]+.")
	match := reMajor.FindStringSubmatch(cliVer)

	cliMajor := match[0][1 : len(match[0])-1]

	pluginInfo := PluginVersions(source)

	availableVersions := pluginInfo[plugin][cliMajor]

	if len(availableVersions) == 0 {
		return ""
	}

	return availableVersions[len(availableVersions)-1]
}
