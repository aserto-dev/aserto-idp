package version

import (
	"fmt"

	"github.com/aserto-dev/idp-plugin-sdk/version"
)

// values set by linker using ldflag -X.
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func getVersion() (string, string, string) {
	return ver, date, commit
}

func GetVer() string {
	buildInfo := version.GetBuildInfo(getVersion)

	return buildInfo.Version
}

func GetVersionString() string {
	buildInfo := version.GetBuildInfo(getVersion)

	return fmt.Sprintf("%s %s %s-%s [%s]",
		buildInfo.Version,
		buildInfo.Commit,
		buildInfo.Os,
		buildInfo.Arch,
		buildInfo.Date,
	)
}
