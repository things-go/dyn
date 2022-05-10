package main

import (
	"github.com/things-go/anyhow/builder"
)

func main() {
	builder.Println()
	builder.Println(builder.WithDeploy("dev"))
	builder.Println(builder.WithMetadata(map[string]string{
		"MetaKey1": "MetaVal1",
		"MetaKey2": "MetaVal2",
	}))
	builder.Println(
		builder.WithDeploy("dev"),
		builder.WithMetadata(map[string]string{
			"MetaKey1": "MetaVal1",
			"MetaKey2": "MetaVal2",
		}),
	)
}
