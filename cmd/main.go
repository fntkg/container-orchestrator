package main

import (
	"github.com/fntkg/container-orchestrator/pkg/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fntkg/container-orchestrator/pkg/controller"
	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
)

func main() {
	// Create a Node Manager and register some nodes.
	nm := node.NewManager()
	if err := nm.Register(node.Node{ID: "node-1", Healthy: true}); err != nil {
		log.Fatalf("Error registering node-1: %v", err)
	}
	if err := nm.Register(node.Node{ID: "node-2", Healthy: true}); err != nil {
		log.Fatalf("Error registering node-2: %v", err)
	}

	// Create some sample tasks.
	tasks := []scheduler.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}

	// Initialize the DefaultScheduler.
	sched := scheduler.NewDefaultScheduler()

	// Create a new Controller Manager with the Node Manager.
	ctrlManager := controller.NewControllerManager(sched, tasks, nm)

	// Run the Controller Manager.
	stopCh := make(chan struct{})
	go ctrlManager.Run(stopCh)

	// Initialize the API router.
	router := api.NewRouter()
	// Start the API server in its own goroutine.
	apiPort := ":8080"
	go func() {
		log.Printf("Starting API server on port %s", apiPort)
		if err := http.ListenAndServe(apiPort, router); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	// Wait for an OS signal (SIGINT, SIGTERM) to gracefully shut down.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	close(stopCh)
	log.Println("Shutting down gracefully...")
}
