package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const indexFile = "index.html"

var (
	port   = flag.Int("p", 8080, "port number")
	folder = flag.String("f", "web", "path to folder")
)

func handle(rw http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(*folder, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(rw, r, filepath.Join(*folder, indexFile))
		return
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(*folder)).ServeHTTP(rw, r)
}

func main() {
	mux := http.NewServeMux()

	// Liveness and readiness probes
	mux.HandleFunc("/live", func(rw http.ResponseWriter, _ *http.Request) { rw.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/ready", func(rw http.ResponseWriter, _ *http.Request) { rw.WriteHeader(http.StatusOK) })

	mux.HandleFunc("/", handle)

	log.Printf("serving %s on port %d", *folder, *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), mux))
}
