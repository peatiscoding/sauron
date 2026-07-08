package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed static
var staticFiles embed.FS

func main() {
	hub := newHub()
	go hub.run()

	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})
	mux.HandleFunc("/api/focus", func(w http.ResponseWriter, r *http.Request) {
		handleFocus(hub, w, r)
	})
	mux.HandleFunc("/api/anchor", func(w http.ResponseWriter, r *http.Request) {
		handleAnchor(hub, w, r)
	})
	mux.HandleFunc("/api/image", handleImage)
	mux.HandleFunc("/", spaHandler(sub))

	log.Println("Sauron eye open on :6905")
	if err := http.ListenAndServe(":6905", mux); err != nil {
		log.Fatal(err)
	}
}

func spaHandler(fsys fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(fsys))
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		f, err := fsys.Open(path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Path not found → SPA fallback: serve index.html
		data, readErr := fs.ReadFile(fsys, "index.html")
		if readErr != nil {
			http.Error(w, "Frontend not built — run scripts/build.sh", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}
}
