// File: pkg/api/router_test.go
package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/api"
	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
)

// FakeNodeManager implements the node.NodeManager interface for testing purposes.
type FakeNodeManager struct {
	nodes []models.Node
}

func (fnm *FakeNodeManager) Register(n models.Node) error {
	// For simplicity, we allow duplicate registration.
	fnm.nodes = append(fnm.nodes, n)
	return nil
}

func (fnm *FakeNodeManager) GetNodes() []models.Node {
	return fnm.nodes
}

func (fnm *FakeNodeManager) UpdateHealth(id string, healthy bool) error {
	for i, n := range fnm.nodes {
		if n.ID == id {
			fnm.nodes[i].Healthy = healthy
			return nil
		}
	}
	return fmt.Errorf("node not found")
}

// Test the /health endpoint.
func TestHealthEndpoint(t *testing.T) {
	fnm := &FakeNodeManager{}
	ds := datastore.NewInMemoryDatastore()
	apiInstance := api.NewAPI(fnm, ds)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("expected body 'OK', got '%s'", string(body))
	}
}

// Test the /nodes GET endpoint.
func TestGetNodesEndpoint(t *testing.T) {
	// Pre-populate the fake NodeManager with some nodes.
	fnm := &FakeNodeManager{
		nodes: []models.Node{
			{ID: "node-1", Healthy: true},
			{ID: "node-2", Healthy: false},
		},
	}
	ds := datastore.NewInMemoryDatastore()
	apiInstance := api.NewAPI(fnm, ds)

	req := httptest.NewRequest("GET", "/nodes", nil)
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var nodes []models.Node
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(nodes))
	}
}

// Test the /nodes POST endpoint.
func TestRegisterNodeEndpoint(t *testing.T) {
	fnm := &FakeNodeManager{}
	ds := datastore.NewInMemoryDatastore()
	apiInstance := api.NewAPI(fnm, ds)

	newNode := models.Node{ID: "node-3", Healthy: true}
	bodyBytes, _ := json.Marshal(newNode)
	req := httptest.NewRequest("POST", "/nodes", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var nodeResp models.Node
	if err := json.NewDecoder(resp.Body).Decode(&nodeResp); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if nodeResp.ID != "node-3" {
		t.Errorf("expected node ID 'node-3', got '%s'", nodeResp.ID)
	}
}

// Test the /nodes/{id} PUT endpoint.
func TestUpdateNodeEndpoint(t *testing.T) {
	fnm := &FakeNodeManager{
		nodes: []models.Node{
			{ID: "node-1", Healthy: true},
		},
	}
	ds := datastore.NewInMemoryDatastore()
	apiInstance := api.NewAPI(fnm, ds)

	payload := map[string]bool{"healthy": false}
	payloadBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/nodes/node-1", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if result["status"] != "updated" {
		t.Errorf("expected status 'updated', got '%s'", result["status"])
	}

	// Verify that the node's health was updated.
	nodes := fnm.GetNodes()
	if len(nodes) != 1 || nodes[0].Healthy != false {
		t.Errorf("expected node health to be false, got %v", nodes[0].Healthy)
	}
}

// Test the /tasks GET endpoint.
func TestGetTasksEndpoint(t *testing.T) {
	// Use the real in-memory datastore and pre-populate with tasks.
	ds := datastore.NewInMemoryDatastore()
	tasks := []models.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}
	for _, task := range tasks {
		if err := ds.SaveTask(task); err != nil {
			t.Fatalf("error saving task: %v", err)
		}
	}

	fnm := &FakeNodeManager{}
	apiInstance := api.NewAPI(fnm, ds)

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var tasksResp []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasksResp); err != nil {
		t.Errorf("error decoding tasks: %v", err)
	}
	if len(tasksResp) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasksResp))
	}
}

// Test the /tasks POST endpoint.
func TestRegisterTaskEndpoint(t *testing.T) {
	ds := datastore.NewInMemoryDatastore()
	fnm := &FakeNodeManager{}
	apiInstance := api.NewAPI(fnm, ds)

	newTask := models.Task{ID: "task-3"}
	bodyBytes, _ := json.Marshal(newTask)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	apiInstance.Router().ServeHTTP(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Errorf("error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var taskResp models.Task
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if taskResp.ID != "task-3" {
		t.Errorf("expected task ID 'task-3', got '%s'", taskResp.ID)
	}
}
