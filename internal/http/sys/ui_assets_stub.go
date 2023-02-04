//go:build !sys_ui

package sys

import "net/http"

func init() {
	uiBuiltIn = false
}

// sysAssetFS is a stub for building Rumors sys without UI.
func sysAssetFS() http.FileSystem {
	return nil
}
