package scheduler

type Task struct {
	ID       string
	Priority int
	Payload  string
}

type Scheduler struct{}

func New() *Scheduler {
	return &Scheduler{}
}
