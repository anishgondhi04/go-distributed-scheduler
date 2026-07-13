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
	fifo     []*Task
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

	if s.strategy == PriorityQ {
		heap.Push(s.queue, t)
	} else {
		s.fifo = append(s.fifo, t)
	}
}
func (s *Scheduler) Next() *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var t *Task

	switch s.strategy {
	case PriorityQ:
		if s.queue.Len() == 0 {
			return nil
		}
		t = heap.Pop(s.queue).(*Task)
	default:
		if len(s.fifo) == 0 {
			return nil
		}
		t = s.fifo[0]
		s.fifo = s.fifo[1:]
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
	if s.strategy == PriorityQ {
		return s.queue.Len()
	}
	return len(s.fifo)
}
func (s *Scheduler) SetStrategy(strategy Strategy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strategy = strategy
}

func (s *Scheduler) GetStrategy() Strategy {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.strategy
}
