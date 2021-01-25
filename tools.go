// +build tools
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package main

import (
	_ "github.com/go-bindata/go-bindata/"
	_ "golang.org/x/tools/cmd/goimports"
)
