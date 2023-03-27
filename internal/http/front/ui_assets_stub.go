//go:build !ui

package front

import "net/http"

func init() {
	uiBuiltIn = false
}

// assetFS is a stub for building Rumors sys without UI.
func assetFS() http.FileSystem {
	return nil
}
