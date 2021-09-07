package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (web *Web) getWorkers(w http.ResponseWriter, r *http.Request) {
	out, err := json.Marshal(web.c.GetWorkers())

	if err != nil {
		w.WriteHeader(500)

		fmt.Fprintf(w, "%v", err)

		return
	}

	fmt.Fprint(w, string(out))
}

func (web *Web) postPackage(w http.ResponseWriter, r *http.Request) {
	var pkgs []string

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&pkgs)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to decode JSON, %v", err)

		return
	}

	log.Printf("packages: %v", pkgs)
}
