// +build dev

package index

//go:generate go run -tags=dev ../../cmd/server/index_generate.go

import "net/http"

// Assets contains project assets.
var Index http.FileSystem = http.Dir("../fs/index")
