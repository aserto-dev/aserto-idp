package finder

import (
	"os"
	"path/filepath"

	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type HomeDir struct {
}

func NewHomeDir() Finder {
	return &HomeDir{}
}

func (path HomeDir) Find() ([]string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	pluginFiles := []string{}

	files, err := filepath.Glob(filepath.Join(homeDir, ".aserto", "idpplugins", x.PluginPrefix+"*"))
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		pluginFiles = append(pluginFiles, files...)
	}

	return pluginFiles, nil
}
