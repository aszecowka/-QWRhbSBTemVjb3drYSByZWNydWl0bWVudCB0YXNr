// +build tools

package weather

// tracking tool dependencies: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/kisielk/errcheck"
	_ "github.com/vektra/mockery"
)
