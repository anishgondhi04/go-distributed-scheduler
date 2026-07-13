package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anishgondhi04/go-distributed-scheduler/internal/node"
	"github.com/anishgondhi04/go-distributed-scheduler/internal/scheduler"
	"github.com/anishgondhi04/go-distributed-scheduler/internal/simulation"
)

func main() {
	nodeMgr := node.NewManager()
	nodeMgr.Register("node-1")
	nodeMgr.Register("node-2")
	nodeMgr.Register("node-3")

	sched := scheduler.New(scheduler.PriorityQ)
	sched.SetNodes(nodeMgr.IDs())

	sim := simulation.New(sched, nodeMgr, 1*time.Second)
	sim.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	addr := ":8080"
	log.Printf("go-distributed-scheduler listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
