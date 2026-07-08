package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all origins — this runs on localhost only.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// focusRequest is the body of POST /api/focus.
type focusRequest struct {
	Content   string `json:"content"`
	Filetype  string `json:"filetype"`
	Filename  string `json:"filename"`
	Filepath  string `json:"filepath"`
	GitStatus string `json:"gitStatus"`
}

func handleFocus(hub *Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req focusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Normalize filetype so the frontend only needs to check two values.
	ft := strings.ToLower(strings.TrimSpace(req.Filetype))
	switch ft {
	case "md":
		ft = "markdown"
	case "yml":
		ft = "yaml"
	}

	hub.publish(Message{
		Type:      "focus",
		Filetype:  ft,
		Filename:  req.Filename,
		Filepath:  req.Filepath,
		Content:   req.Content,
		GitStatus: req.GitStatus,
	})

	log.Printf("focus: %s (%s)", req.Filename, ft)
	w.WriteHeader(http.StatusNoContent)
}

// anchorRequest is the body of POST /api/anchor.
type anchorRequest struct {
	Line  int `json:"line"`
	Total int `json:"total"`
}

func handleAnchor(hub *Hub, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req anchorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Total <= 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	hub.publishNoStore(Message{Type: "anchor", Line: req.Line, Total: req.Total})
	w.WriteHeader(http.StatusNoContent)
}

var imageTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".avif": "image/avif",
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	if p == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}
	p = filepath.Clean(p)
	mime, ok := imageTypes[strings.ToLower(filepath.Ext(p))]
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	data, err := os.ReadFile(p)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", mime)
	w.Write(data)
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}

	client := &WSClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 32),
	}
	hub.register <- client

	go client.writePump()
	go client.readPump()
}
