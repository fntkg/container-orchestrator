package datastore

import (
	"github.com/fntkg/container-orchestrator/pkg/models"
	"testing"
)

func TestInMemoryDatastore_SaveAndGetNodes(t *testing.T) {
	// Create a new in-memory datastore.
	ds := NewInMemoryDatastore()

	// Define some dummy nodes.
	dummyNodes := []models.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: false},
	}

	// Save nodes in the datastore.
	for _, n := range dummyNodes {
		if err := ds.SaveNode(n); err != nil {
			t.Fatalf("Error saving node %s: %v", n.ID, err)
		}
	}

	// Retrieve nodes.
	nodes, err := ds.GetNodes()
	if err != nil {
		t.Fatalf("Error retrieving nodes: %v", err)
	}

	// Verify that the number of nodes retrieved matches the number of dummy nodes.
	if len(nodes) != len(dummyNodes) {
		t.Errorf("Expected %d nodes, got %d", len(dummyNodes), len(nodes))
	}
}

func TestInMemoryDatastore_SaveAndGetTasks(t *testing.T) {
	// Create a new in-memory datastore.
	ds := NewInMemoryDatastore()

	// Define some dummy tasks.
	dummyTasks := []models.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}

	// Save tasks in the datastore.
	for _, task := range dummyTasks {
		if err := ds.SaveTask(task); err != nil {
			t.Fatalf("Error saving task %s: %v", task.ID, err)
		}
	}

	// Retrieve tasks.
	tasks, err := ds.GetTasks()
	if err != nil {
		t.Fatalf("Error retrieving tasks: %v", err)
	}

	// Verify that the number of tasks retrieved matches the number of dummy tasks.
	if len(tasks) != len(dummyTasks) {
		t.Errorf("Expected %d tasks, got %d", len(dummyTasks), len(tasks))
	}
}
