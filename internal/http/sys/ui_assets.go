//go:build sys_ui

package sys

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/*
var sysContent embed.FS

// sysAssetFS is a http Filesystem that serves the generated sys UI
func sysAssetFS() http.FileSystem {
	f, err := fs.Sub(sysContent, "ui")
	if err != nil {
		panic(err)
	}
	return http.FS(f)
}
