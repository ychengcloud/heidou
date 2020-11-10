// +build ignore

package main

import (
	"log"

	"github.com/decker502/heidou/assets"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(assets.Project, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Project",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
