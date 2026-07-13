package simulation

import (
	"testing"
	"time"

	"github.com/anishgondhi04/go-distributed-scheduler/internal/node"
	"github.com/anishgondhi04/go-distributed-scheduler/internal/scheduler"
)

func TestSimulatorDispatches(t *testing.T) {
	nodeMgr := node.NewManager()
	nodeMgr.Register("node-1")

	sched := scheduler.New(scheduler.PriorityQ)
	sched.SetNodes(nodeMgr.IDs())

	sim := New(sched, nodeMgr, 50*time.Millisecond)
	sim.Start()

	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	n, ok := nodeMgr.Get("node-1")
	if !ok {
		t.Fatal("expected node-1 to exist")
	}
	if n.LastActiveAt.IsZero() {
		t.Fatal("expected node to have been active")
	}
}
