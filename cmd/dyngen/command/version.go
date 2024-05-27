package command

import (
	"fmt"
	"runtime"
)

const version = "v1.0.0"

func buildVersion() string {
	return fmt.Sprintf("%s\nGo Version: %s\nGo Os: %s\nGo Arch: %s\n",
		version, runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
}
