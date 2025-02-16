package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fntkg/container-orchestrator/pkg/api"
	"github.com/fntkg/container-orchestrator/pkg/controller"
	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
	"github.com/fntkg/container-orchestrator/pkg/taskmanager"
)

func main() {
	// Initialize the in-memory datastore.
	ds := datastore.NewInMemoryDatastore()

	// Create the Node DefaultNodeManager using the datastore.
	nm := node.NewManager(ds)
	if err := nm.Register(models.Node{ID: "node-1", Healthy: true}); err != nil {
		log.Fatalf("Failed to register node-1: %v", err)
	}
	if err := nm.Register(models.Node{ID: "node-2", Healthy: true}); err != nil {
		log.Fatalf("Failed to register node-2: %v", err)
	}

	// Create the Task DefaultNodeManager using the datastore.
	tm := taskmanager.NewTaskManager(ds)
	// Optionally, create some initial tasks.
	if err := tm.CreateTask(models.Task{ID: "task-1", Status: "pending"}); err != nil {
		log.Fatalf("Failed to create task-1: %v", err)
	}
	if err := tm.CreateTask(models.Task{ID: "task-2", Status: "pending"}); err != nil {
		log.Fatalf("Failed to create task-2: %v", err)
	}

	// Initialize the scheduler.
	sched := scheduler.NewDefaultScheduler()

	// Create the Controller DefaultNodeManager with the scheduler, Task DefaultNodeManager, and Node DefaultNodeManager.
	ctrlManager := controller.NewControllerManager(sched, tm, nm)
	stopCh := make(chan struct{})
	go ctrlManager.Run(stopCh)

	// Create the API router with the Node DefaultNodeManager and datastore.
	apiInstance := api.NewAPI(nm, tm)
	apiPort := ":8080"
	go func() {
		log.Printf("Starting API server on port %s", apiPort)
		if err := http.ListenAndServe(apiPort, apiInstance.Router()); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	// Listen for OS signals to gracefully shut down.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down gracefully...")
	close(stopCh)
}
