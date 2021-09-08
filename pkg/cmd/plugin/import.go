package plugin

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type ImportCmd struct {
	InclUserExt bool
	Source      string
	kong.Plugins
}

func (cmd *ImportCmd) Run(c *cc.CC) error {
	return nil
}
