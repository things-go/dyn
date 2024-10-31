package meta

import (
	"fmt"
	"runtime"
)

const Version = "v1.0.0"

func BuildVersion() string {
	return fmt.Sprintf("%s\nGo Version: %s\nGo Os: %s\nGo Arch: %s\n",
		Version, runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
}
