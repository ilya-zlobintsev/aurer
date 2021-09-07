//go:build prod

package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed public
var embedFrontend embed.FS

func getFrontendAssets() fs.FS {
	log.Println("Serving embedded web assets")

	f, err := fs.Sub(embedFrontend, "public")
	if err != nil {
		panic(err)
	}

	return f
}
