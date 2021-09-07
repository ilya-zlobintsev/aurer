package web

import (
	"io/fs"
	"log"
	"net/http"

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
	r.HandleFunc("/api/packages", web.postPackage).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.FS(web.frontend)))

	http.ListenAndServe(":8008", logRequest(r))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
