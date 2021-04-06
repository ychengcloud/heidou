package assets

import (
	"embed"
)

//go:embed doc/*
var Swagger embed.FS
