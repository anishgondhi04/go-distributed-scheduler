package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/anishgondhi04/go-distributed-scheduler/internal/node"
	"github.com/anishgondhi04/go-distributed-scheduler/internal/scheduler"
)

type Simulator struct {
	sched    *scheduler.Scheduler
	nodes    *node.Manager
	interval time.Duration
	stopCh   chan struct{}
}

func New(sched *scheduler.Scheduler, nodes *node.Manager, interval time.Duration) *Simulator {
	return &Simulator{
		sched:    sched,
		nodes:    nodes,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}
func (s *Simulator) Start() {
	go s.run()
}

func (s *Simulator) Stop() {
	close(s.stopCh)
}

func (s *Simulator) run() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	taskNum := 0

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			taskNum++
			t := &scheduler.Task{
				ID:       fmt.Sprintf("task-%d", taskNum),
				Priority: rand.Intn(10) + 1,
				Payload:  fmt.Sprintf("payload-%d", taskNum),
			}
			s.sched.Submit(t)
			s.dispatch()
		}
	}
}

func (s *Simulator) dispatch() {
	next := s.sched.Next()
	if next == nil {
		return
	}

	nodeID := s.nodes.LeastLoaded()
	if nodeID == "" {
		return
	}

	next.NodeID = nodeID
	s.nodes.AssignTask(nodeID)

	go func() {
		time.Sleep(500 * time.Millisecond)
		s.nodes.CompleteTask(nodeID)
	}()
	//log.Printf("dispatched %s to %s", next.ID, nodeID)
}
