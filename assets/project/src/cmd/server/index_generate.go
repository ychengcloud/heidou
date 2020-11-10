// +build ignore

package main

import (
	"log"

	"github.com/h/d/graphql/index"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(index.Index, vfsgen.Options{
		PackageName:  "index",
		BuildTags:    "!dev",
		VariableName: "Index",
	})
	if err != nil {
		log.Fatalln(err)
	}

}
