package api

import (
	"encoding/json"
	"net/http"

	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/gorilla/mux"
)

// API encapsulates the HTTP router and its dependencies.
type API struct {
	router      *mux.Router
	nodeManager node.NodeManager // NodeManager interface
	ds          datastore.Datastore
}

// NewAPI creates a new API instance with the provided NodeManager and Datastore.
func NewAPI(nm node.NodeManager, ds datastore.Datastore) *API {
	r := mux.NewRouter().StrictSlash(true)
	api := &API{
		router:      r,
		nodeManager: nm,
		ds:          ds,
	}

	// Health endpoint
	r.HandleFunc("/health", api.healthHandler).Methods("GET")

	// Node endpoints
	r.HandleFunc("/nodes", api.getNodesHandler).Methods("GET")
	r.HandleFunc("/nodes", api.registerNodeHandler).Methods("POST")
	r.HandleFunc("/nodes/{id}", api.updateNodeHandler).Methods("PUT")

	// Task endpoints
	r.HandleFunc("/tasks", api.getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", api.registerTaskHandler).Methods("POST")

	return api
}

// Router returns the underlying mux.Router.
func (a *API) Router() *mux.Router {
	return a.router
}

// healthHandler returns a simple "OK" to indicate the service is up.
func (a *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// getNodesHandler returns the list of registered nodes.
func (a *API) getNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := a.nodeManager.GetNodes()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(nodes)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// registerNodeHandler registers a new node.
func (a *API) registerNodeHandler(w http.ResponseWriter, r *http.Request) {
	var n models.Node
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := a.nodeManager.Register(n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(n)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// updateNodeHandler updates the health status of an existing node.
func (a *API) updateNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Expecting a JSON payload with a "healthy" field.
	var payload struct {
		Healthy bool `json:"healthy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := a.nodeManager.UpdateHealth(id, payload.Healthy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// getTasksHandler returns the list of registered tasks.
func (a *API) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := a.ds.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// registerTaskHandler registers a new task.
func (a *API) registerTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := a.ds.SaveTask(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
