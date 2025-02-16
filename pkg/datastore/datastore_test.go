// File: pkg/datastore/datastore_test.go
package datastore_test

import (
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
)

func TestInMemoryDatastore_SaveAndGetNodes(t *testing.T) {
	// Create a new in-memory datastore.
	ds := datastore.NewInMemoryDatastore()

	// Define some sample nodes.
	nodesToSave := []models.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: false},
	}

	// Save each node.
	for _, n := range nodesToSave {
		if err := ds.SaveNode(n); err != nil {
			t.Fatalf("Failed to save node %s: %v", n.ID, err)
		}
	}

	// Retrieve nodes from the datastore.
	retrievedNodes, err := ds.GetNodes()
	if err != nil {
		t.Fatalf("Error retrieving nodes: %v", err)
	}

	// Verify the number of nodes retrieved.
	if len(retrievedNodes) != len(nodesToSave) {
		t.Errorf("Expected %d nodes, got %d", len(nodesToSave), len(retrievedNodes))
	}

	// Optionally, verify that each saved node exists in the retrieved slice.
	nodeMap := make(map[string]models.Node)
	for _, n := range retrievedNodes {
		nodeMap[n.ID] = n
	}
	for _, n := range nodesToSave {
		if savedNode, ok := nodeMap[n.ID]; !ok {
			t.Errorf("Node %s was not found", n.ID)
		} else if savedNode.Healthy != n.Healthy {
			t.Errorf("Node %s: expected Healthy %v, got %v", n.ID, n.Healthy, savedNode.Healthy)
		}
	}
}

func TestInMemoryDatastore_SaveAndGetTasks(t *testing.T) {
	// Create a new in-memory datastore.
	ds := datastore.NewInMemoryDatastore()

	// Define some sample tasks.
	tasksToSave := []models.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}

	// Save each task.
	for _, task := range tasksToSave {
		if err := ds.SaveTask(task); err != nil {
			t.Fatalf("Failed to save task %s: %v", task.ID, err)
		}
	}

	// Retrieve tasks from the datastore.
	retrievedTasks, err := ds.GetTasks()
	if err != nil {
		t.Fatalf("Error retrieving tasks: %v", err)
	}

	// Verify the number of tasks retrieved.
	if len(retrievedTasks) != len(tasksToSave) {
		t.Errorf("Expected %d tasks, got %d", len(tasksToSave), len(retrievedTasks))
	}

	// Optionally, verify that each saved task exists in the retrieved slice.
	taskMap := make(map[string]models.Task)
	for _, tsk := range retrievedTasks {
		taskMap[tsk.ID] = tsk
	}
	for _, tsk := range tasksToSave {
		if _, ok := taskMap[tsk.ID]; !ok {
			t.Errorf("Task %s was not found", tsk.ID)
		}
	}
}
