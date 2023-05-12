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
	port = flag.Int("port", 1234, "port number")
	path = flag.String("path", "spa", "path to SPA folder")
)

type handler struct {
	path string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		log.Printf("failed to parse path %s: %v", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.path, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("path %s does not exist", path)
		http.ServeFile(w, r, filepath.Join(h.path, indexFile))
		return
	}
	if err != nil {
		log.Printf("failed to get file info for path %s: %v", path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.path)).ServeHTTP(w, r)
}

func main() {
	http.Handle("/", handler{path: *path})

	log.Printf("serving %s on port %d", *path, *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
