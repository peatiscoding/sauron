package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message is the wire protocol between server and browser.
type Message struct {
	Type      string `json:"type"`
	Filetype  string `json:"filetype,omitempty"`
	Filename  string `json:"filename,omitempty"`
	Filepath  string `json:"filepath,omitempty"`
	Content   string `json:"content,omitempty"`
	GitStatus string `json:"gitStatus,omitempty"`
	Line      int    `json:"line,omitempty"`
	Total     int    `json:"total,omitempty"`
}

// WSClient is a connected browser tab.
type WSClient struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub maintains the set of active WebSocket clients and broadcasts to them.
type Hub struct {
	mu         sync.Mutex
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	last       []byte // replayed to new connections immediately
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan []byte, 32),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			last := h.last
			h.mu.Unlock()
			// Send current state so the new tab sees content immediately.
			if last != nil {
				select {
				case c.send <- last:
				default:
				}
			}

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			h.last = msg
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) publish(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	h.broadcast <- data
}

// publishNoStore broadcasts msg to all clients without updating h.last.
// Use for ephemeral messages (e.g. anchor) that should not be replayed to new tabs.
func (h *Hub) publishNoStore(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	h.mu.Lock()
	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
	h.mu.Unlock()
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *WSClient) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readPump reads (and discards) incoming frames to detect disconnects.
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}
