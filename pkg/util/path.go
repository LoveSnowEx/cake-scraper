package util

import (
	"path"
	"runtime"
)

var ProjectRoot string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	ProjectRoot = path.Join(path.Dir(filename), "../..")
}
