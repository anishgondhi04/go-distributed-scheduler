package scheduler

import "testing"

func TestSchedulerBasic(t *testing.T) {
	s := New(PriorityQ)
	s.SetNodes([]string{"node-1", "node-2"})

	s.Submit(&Task{ID: "t1", Priority: 1})
	s.Submit(&Task{ID: "t2", Priority: 5})

	first := s.Next()
	if first.ID != "t2" {
		t.Fatalf("expected t2 (higher priority) first, got %s", first.ID)
	}
}
