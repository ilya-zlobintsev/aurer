package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ilyazzz/aurer/internal/repo"
)

func respondJson(w http.ResponseWriter, j interface{}) {
	if err := json.NewEncoder(w).Encode(j); err != nil {
		w.WriteHeader(500)

		fmt.Fprintf(w, "%v", err)
	}
}

func (web *Web) getWorkers(w http.ResponseWriter, r *http.Request) {
	workers := web.c.GetWorkers()

	respondJson(w, workers)
}

func (web *Web) getPackages(w http.ResponseWriter, r *http.Request) {
	packages, err := repo.ReadRepo(web.c.RepoDir)

	if err != nil {
		w.WriteHeader(500)

		fmt.Fprintf(w, "%v", err)

		return
	}

	respondJson(w, packages)
}

func (web *Web) postPackage(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)

		log.Printf("Failed reading request body: %v", b)

		return
	}

	pkgname := string(b)

	go func() {
		err := web.c.BuildPackage(pkgname)

		if err != nil {
			log.Printf("Error building package %v: %v", pkgname, err)
		}
	}()
}
