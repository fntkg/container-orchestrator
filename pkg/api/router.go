package api

import (
	"encoding/json"
	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
	"github.com/gorilla/mux" // You can install this package with: go get -u github.com/gorilla/mux
	"log"
	"net/http"
)

// NewRouter configures and returns an HTTP router.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Health endpoint
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// Schedule endpoint.
	router.HandleFunc("/schedule", ScheduleHandler).Methods("POST")
	return router
}

// HealthCheckHandler is a simple endpoint to verify that the service is running.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		// The error is logged if fail to write the answer.
		log.Printf("Error writing the response: %v", err)
	}
}

// ScheduleHandler receives a task, calls the scheduler, and returns the assigned node.
func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var task scheduler.Task
	// Decode the incoming JSON payload into a Task struct.
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task payload", http.StatusBadRequest)
		return
	}

	// TODO: Create real nodes
	// Simulate available nodes.
	nodes := []node.Node{
		{ID: "node-1"},
		{ID: "node-2"},
	}

	// Create an instance of the DefaultScheduler.
	sched := scheduler.NewDefaultScheduler()
	selectedNode, err := sched.Schedule(task, nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the assigned node as JSON.
	response := map[string]string{"assigned_node": selectedNode.ID}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
