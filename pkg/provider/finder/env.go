package finder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto-idp/pkg/x"
)

type Environment struct {
}

func NewEnvironment() Finder {
	return &Environment{}
}

func (path Environment) Find() ([]string, error) {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dirs = append(dirs, pwd)

	pluginFiles := []string{}
	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, x.PluginPrefix+"*"))
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			pluginFiles = append(pluginFiles, files...)
		}
	}

	return pluginFiles, nil
}
