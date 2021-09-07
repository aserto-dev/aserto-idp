package plugin

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto-idp/pkg/cc"
)

type ImportCmd struct {
	InclUserExt bool
	Source      string
	context     *cc.CC
	kong.Plugins
}

func (cmd *ImportCmd) Run(c *cc.CC) error {
	fmt.Println("xxxx" + c.Provider + "xxx")
	return nil
}
