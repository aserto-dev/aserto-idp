package cmd //nolint:testpackage // dedicated package to export private type to be used in _test namespace.

// export private package functions to tests.
var (
	DownloadProvider = downloadProvider
	CheckForUpdates  = checkForUpdates
	ValidatePlugin   = validatePlugin
)
