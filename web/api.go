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

func handleErr(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, err)
}

func bodyStr(w http.ResponseWriter, r *http.Request) (string, error) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (web *Web) getWorkers(w http.ResponseWriter, r *http.Request) {
	workers := web.c.GetWorkers()

	respondJson(w, workers)
}

func (web *Web) getPackages(w http.ResponseWriter, r *http.Request) {
	packages, err := repo.ReadRepo(web.c.RepoDir)

	if len(packages) == 0 || err != nil {
		respondJson(w, make([]repo.PkgInfo, 0))

		return
	}

	respondJson(w, packages)
}
func (web *Web) postPackage(w http.ResponseWriter, r *http.Request) {
	pkgname, err := bodyStr(w, r)

	if err != nil {
		handleErr(w, err)

		return
	}

	go func() {
		err := web.c.BuildPackage(pkgname)

		if err != nil {
			log.Printf("Error building package %v: %v", pkgname, err)
		}
	}()
}

func (web *Web) startUpdate(w http.ResponseWriter, r *http.Request) {
	err := web.c.Update()

	if err != nil {
		handleErr(w, err)
	}
}

func (web *Web) deletePackage(w http.ResponseWriter, r *http.Request) {
	pkgName, err := bodyStr(w, r)

	if err != nil {
		handleErr(w, err)

		return
	}

	err = repo.DeletePackage(web.c.RepoDir, pkgName)

	if err != nil {
		handleErr(w, err)

	}
}
