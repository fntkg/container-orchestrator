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
)

func main() {
	// Initialize the in-memory datastore.
	ds := datastore.NewInMemoryDatastore()

	// Create the Node Manager using the datastore.
	nm := node.NewManager(ds)

	// Register some nodes via the Node Manager.
	if err := nm.Register(models.Node{ID: "node-1", Healthy: true}); err != nil {
		log.Fatalf("Failed to register node-1: %v", err)
	}
	if err := nm.Register(models.Node{ID: "node-2", Healthy: true}); err != nil {
		log.Fatalf("Failed to register node-2: %v", err)
	}

	// Save some sample tasks in the datastore.
	tasks := []models.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}
	for _, task := range tasks {
		if err := ds.SaveTask(models.Task{ID: task.ID}); err != nil {
			log.Fatalf("Failed to save task %s: %v", task.ID, err)
		}
	}

	// Initialize the scheduler and Controller Manager.
	sched := scheduler.NewDefaultScheduler()
	ctrlManager := controller.NewControllerManager(sched, tasks, nm)
	stopCh := make(chan struct{})
	go ctrlManager.Run(stopCh)

	// Create the API router with the Node Manager and Datastore.
	apiInstance := api.NewAPI(nm, ds)

	// Start the HTTP server.
	apiPort := ":8080"
	go func() {
		log.Printf("Starting API server on port %s", apiPort)
		if err := http.ListenAndServe(apiPort, apiInstance.Router()); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()

	// Wait for OS signals to gracefully shut down.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down gracefully...")
	close(stopCh)
}
