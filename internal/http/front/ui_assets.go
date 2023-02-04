//go:build ui

package front

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/*
var content embed.FS

// assetFS is a http Filesystem that serves the generated sys UI
func assetFS() http.FileSystem {
	f, err := fs.Sub(content, "ui")
	if err != nil {
		panic(err)
	}
	return http.FS(f)
}
