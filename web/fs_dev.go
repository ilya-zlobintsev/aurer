//go:build !prod

package web

import (
	"io/fs"
	"log"
	"os"
)

func getFrontendAssets() fs.FS {
	log.Println("Serving web assets from public/")

	return os.DirFS("../../web/frontend/public")
}
