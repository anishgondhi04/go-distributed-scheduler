package node

import "testing"

func TestManagerLeastLoaded(t *testing.T) {
	m := NewManager()
	m.Register("node-1")
	m.Register("node-2")

	m.AssignTask("node-1")
	m.AssignTask("node-1")
	m.AssignTask("node-2")

	if got := m.LeastLoaded(); got != "node-2" {
		t.Fatalf("expected node-2 to be least loaded, got %s", got)
	}
}
