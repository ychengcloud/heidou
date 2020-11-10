// +build dev

package assets

//go:generate go run -tags=dev ../../cmd/server/assets_generate.go

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("../fs/assets")
