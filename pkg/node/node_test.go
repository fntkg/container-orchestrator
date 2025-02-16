package node

import (
	"fmt"
	"testing"

	_ "github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
)

// FakeDatastore is a mock implementation of datastore.Datastore for testing purposes.
type FakeDatastore struct {
	nodes map[string]models.Node
	// We include tasks as no-op to satisfy the interface.
	tasks map[string]models.Task
}

func NewFakeDatastore() *FakeDatastore {
	return &FakeDatastore{
		nodes: make(map[string]models.Node),
		tasks: make(map[string]models.Task),
	}
}

// SaveNode saves a node. It returns an error if a node with the same ID and identical Healthy status already exists.
// If the Healthy status is different (indicating an update), it updates the stored node.
func (fds *FakeDatastore) SaveNode(n models.Node) error {
	if existing, ok := fds.nodes[n.ID]; ok {
		if existing.Healthy == n.Healthy {
			return fmt.Errorf("node %s already registered", n.ID)
		}
		// Otherwise, update the node (simulate an update operation).
		fds.nodes[n.ID] = n
		return nil
	}
	fds.nodes[n.ID] = n
	return nil
}

// GetNodes retrieves all stored nodes.
func (fds *FakeDatastore) GetNodes() ([]models.Node, error) {
	nodes := make([]models.Node, 0, len(fds.nodes))
	for _, n := range fds.nodes {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

// SaveTask is a stub to satisfy the datastore.Datastore interface.
func (fds *FakeDatastore) SaveTask(t models.Task) error {
	fds.tasks[t.ID] = t
	return nil
}

// GetTasks is a stub to satisfy the datastore.Datastore interface.
func (fds *FakeDatastore) GetTasks() ([]models.Task, error) {
	tasks := make([]models.Task, 0, len(fds.tasks))
	for _, t := range fds.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

//
// Now, the tests for the Node Manager using the FakeDatastore
//

func TestNodeManager_RegisterAndGetNodes(t *testing.T) {
	// Create a Node Manager using the FakeDatastore.
	manager := NewManager(NewFakeDatastore())

	// Define a couple of nodes.
	nodesToRegister := []models.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: true},
	}

	// Register the nodes.
	for _, n := range nodesToRegister {
		if err := manager.Register(n); err != nil {
			t.Errorf("unexpected error registering node %s: %v", n.ID, err)
		}
	}

	// Retrieve nodes and verify the count.
	nodes := manager.GetNodes()
	if len(nodes) != len(nodesToRegister) {
		t.Errorf("expected %d nodes, got %d", len(nodesToRegister), len(nodes))
	}
}

func TestNodeManager_UpdateHealth(t *testing.T) {
	// Create a Node Manager using the FakeDatastore.
	manager := NewManager(NewFakeDatastore())

	// Register a node.
	node := models.Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Update health status (changing Healthy from true to false).
	if err := manager.UpdateHealth("node-1", false); err != nil {
		t.Fatalf("failed to update node health: %v", err)
	}

	// Verify the update.
	nodes := manager.GetNodes()
	var found bool
	for _, n := range nodes {
		if n.ID == "node-1" {
			found = true
			if n.Healthy != false {
				t.Errorf("expected node-1 to be unhealthy, but got healthy")
			}
		}
	}
	if !found {
		t.Errorf("node-1 not found after update")
	}
}

func TestNodeManager_RegisterDuplicate(t *testing.T) {
	// Create a Node Manager using the FakeDatastore.
	manager := NewManager(NewFakeDatastore())

	node := models.Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Attempt to register the same node again (with the same Healthy value).
	err := manager.Register(node)
	if err == nil {
		t.Errorf("expected error when registering duplicate node, got nil")
	}
}
