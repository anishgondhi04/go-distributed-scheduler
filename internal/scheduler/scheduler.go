package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

type Task struct {
	ID        string
	Priority  int
	Payload   string
	CreatedAt time.Time
	NodeID    string
}

type Strategy string

const (
	FCFS       Strategy = "fcfs"
	PriorityQ  Strategy = "priority"
	RoundRobin Strategy = "round_robin"
)

type Scheduler struct {
	mu       sync.Mutex
	strategy Strategy
	queue    *taskHeap
	rrIndex  int
	nodeIDs  []string
}

func New(strategy Strategy) *Scheduler {
	q := &taskHeap{}
	heap.Init(q)
	return &Scheduler{
		strategy: strategy,
		queue:    q,
	}
}
func (s *Scheduler) SetNodes(nodeIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeIDs = nodeIDs
}

func (s *Scheduler) Submit(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t.CreatedAt = time.Now()
	heap.Push(s.queue, t)
}

func (s *Scheduler) Next() *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.queue.Len() == 0 {
		return nil
	}

	var t *Task

	switch s.strategy {
	case PriorityQ:
		t = heap.Pop(s.queue).(*Task)
	default:
		t = heap.Pop(s.queue).(*Task)
	}

	t.NodeID = s.assignNode()
	return t
}

func (s *Scheduler) assignNode() string {
	if len(s.nodeIDs) == 0 {
		return ""
	}
	node := s.nodeIDs[s.rrIndex%len(s.nodeIDs)]
	s.rrIndex++
	return node
}

func (s *Scheduler) QueueLength() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Len()
}
