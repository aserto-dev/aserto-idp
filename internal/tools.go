//go:build tools
// +build tools

package internal

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/wire/cmd/wire"
)
