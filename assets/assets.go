// +build dev

package assets

//go:generate go run -tags=dev ../cmd/heidou/assets_generate.go

import "net/http"

// Assets contains project assets.
var Project http.FileSystem = http.Dir("project")
