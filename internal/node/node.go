package node

import (
	"sync"
	"time"
)

type Status string

const (
	StatusIdle    Status = "idle"
	StatusBusy    Status = "busy"
	StatusOffline Status = "offline"
)

type Node struct {
	ID           string
	Status       Status
	TaskCount    int
	LastActiveAt time.Time
}

type Manager struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

func NewManager() *Manager {
	return &Manager{
		nodes: make(map[string]*Node),
	}
}
func (m *Manager) Register(id string) *Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	n := &Node{
		ID:           id,
		Status:       StatusIdle,
		LastActiveAt: time.Now(),
	}
	m.nodes[id] = n
	return n
}

func (m *Manager) IDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.nodes))
	for id := range m.nodes {
		ids = append(ids, id)
	}
	return ids
}

func (m *Manager) Get(id string) (*Node, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	n, ok := m.nodes[id]
	return n, ok
}

func (m *Manager) AssignTask(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n, ok := m.nodes[id]; ok {
		n.TaskCount++
		n.Status = StatusBusy
		n.LastActiveAt = time.Now()
	}
}

func (m *Manager) CompleteTask(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n, ok := m.nodes[id]; ok {
		if n.TaskCount > 0 {
			n.TaskCount--
		}
		if n.TaskCount == 0 {
			n.Status = StatusIdle
		}
		n.LastActiveAt = time.Now()
	}
}

func (m *Manager) LeastLoaded() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var best *Node
	for _, n := range m.nodes {
		if n.Status == StatusOffline {
			continue
		}
		if best == nil || n.TaskCount < best.TaskCount {
			best = n
		}
	}
	if best == nil {
		return ""
	}
	return best.ID
}

func (m *Manager) Snapshot() []*Node {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]*Node, 0, len(m.nodes))
	for _, n := range m.nodes {
		out = append(out, n)
	}
	return out
}
