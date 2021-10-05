package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ilyazzz/aurer/internal"
)

type Web struct {
	c        *internal.Controller
	frontend fs.FS
}

func InitWeb(c *internal.Controller) Web {

	frontend := getFrontendAssets()

	return Web{
		c:        c,
		frontend: frontend,
	}
}

func (web *Web) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/api/workers", web.getWorkers)
	r.HandleFunc("/api/packages", web.getPackages).Methods("GET")
	r.HandleFunc("/api/packages", web.postPackage).Methods("POST")

	repoServer := http.FileServer(http.Dir(web.c.RepoDir))

	r.PathPrefix("/repo/").Handler(http.StripPrefix("/repo", repoServer))
	r.PathPrefix("/").Handler(http.FileServer(http.FS(web.frontend)))

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":8008", loggedRouter)
}
