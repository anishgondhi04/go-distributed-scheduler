package api

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/anishgondhi04/go-distributed-scheduler/internal/node"
	"github.com/anishgondhi04/go-distributed-scheduler/internal/scheduler"
)

type Metrics struct {
	TasksDispatched int64
	StartedAt       time.Time
}

type Server struct {
	sched   *scheduler.Scheduler
	nodes   *node.Manager
	metrics *Metrics
}

func NewServer(sched *scheduler.Scheduler, nodes *node.Manager) *Server {
	return &Server{
		sched:   sched,
		nodes:   nodes,
		metrics: &Metrics{StartedAt: time.Now()},
	}
}

func (s *Server) IncrementDispatched() {
	atomic.AddInt64(&s.metrics.TasksDispatched, 1)
}
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/nodes", s.handleNodes)
	mux.HandleFunc("/api/queue", s.handleQueue)
	mux.HandleFunc("/api/metrics", s.handleMetrics)
}

func (s *Server) handleNodes(w http.ResponseWriter, r *http.Request) {
	nodes := s.nodes.Snapshot()
	writeJSON(w, nodes)
}

func (s *Server) handleQueue(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]int{
		"queue_length": s.sched.QueueLength(),
	})
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.metrics.StartedAt).Seconds()
	writeJSON(w, map[string]any{
		"tasks_dispatched": atomic.LoadInt64(&s.metrics.TasksDispatched),
		"uptime_seconds":   uptime,
		"queue_length":     s.sched.QueueLength(),
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v)
}
