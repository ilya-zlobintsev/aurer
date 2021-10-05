//go:build prod

package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed frontend/public
var embedFrontend embed.FS

func getFrontendAssets() fs.FS {
	log.Println("Serving embedded web assets")

	f, err := fs.Sub(embedFrontend, "frontend/public")
	if err != nil {
		panic(err)
	}

	return f
}
