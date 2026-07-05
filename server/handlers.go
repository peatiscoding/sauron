package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all origins — this runs on localhost only.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// focusRequest is the body of POST /api/focus.
type focusRequest struct {
	Content  string `json:"content"`
	Filetype string `json:"filetype"`
	Filename string `json:"filename"`
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
		Type:     "focus",
		Filetype: ft,
		Filename: req.Filename,
		Content:  req.Content,
	})

	log.Printf("focus: %s (%s)", req.Filename, ft)
	w.WriteHeader(http.StatusNoContent)
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
