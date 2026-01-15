package assets

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var content embed.FS

func GetFS() (fs.FS, error) {
	return fs.Sub(content, "dist")
}
