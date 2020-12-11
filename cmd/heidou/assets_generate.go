// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"

	"github.com/horcus/heidou/assets"
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
