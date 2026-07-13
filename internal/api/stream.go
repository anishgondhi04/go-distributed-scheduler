package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Broadcaster struct {
	mu      sync.Mutex
	clients map[chan []byte]bool
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan []byte]bool),
	}
}

func (b *Broadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 8)
	b.mu.Lock()
	b.clients[ch] = true
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *Broadcaster) Publish(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- data:
		default:
		}
	}
}
func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := s.broadcaster.Subscribe()
	defer s.broadcaster.Unsubscribe(ch)

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
