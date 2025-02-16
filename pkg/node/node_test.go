// File: pkg/node/node_test.go
package node

import (
	"fmt"
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/models"
)

// FakeDatastore is a fake implementation of datastore.Datastore for testing Node DefaultNodeManager.
type FakeDatastore struct {
	nodes map[string]models.Node
	// Stubs for tasks to satisfy the interface.
	tasks map[string]models.Task
}

// NewFakeDatastore creates a new fake datastore.
func NewFakeDatastore() *FakeDatastore {
	return &FakeDatastore{
		nodes: make(map[string]models.Node),
		tasks: make(map[string]models.Task),
	}
}

// SaveNode simulates saving a node. If the node already exists with the same Healthy value, it returns an error.
func (fds *FakeDatastore) SaveNode(n models.Node) error {
	if existing, ok := fds.nodes[n.ID]; ok {
		if existing.Healthy == n.Healthy {
			return fmt.Errorf("node %s already registered", n.ID)
		}
	}
	fds.nodes[n.ID] = n
	return nil
}

// GetNodes returns all stored nodes.
func (fds *FakeDatastore) GetNodes() ([]models.Node, error) {
	nodes := make([]models.Node, 0, len(fds.nodes))
	for _, n := range fds.nodes {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

// SaveTask and GetTasks are stub methods to satisfy the datastore.Datastore interface.
func (fds *FakeDatastore) SaveTask(t models.Task) error {
	fds.tasks[t.ID] = t
	return nil
}

func (fds *FakeDatastore) GetTasks() ([]models.Task, error) {
	tasks := make([]models.Task, 0, len(fds.tasks))
	for _, t := range fds.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func TestNodeManager_RegisterAndGetNodes(t *testing.T) {
	ds := NewFakeDatastore()
	manager := NewManager(ds)

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
	ds := NewFakeDatastore()
	manager := NewManager(ds)

	// Register a node.
	node1 := models.Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node1); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Update its health status.
	if err := manager.UpdateHealth("node-1", false); err != nil {
		t.Fatalf("failed to update node health: %v", err)
	}

	// Verify that the node's health was updated.
	nodes := manager.GetNodes()
	var updated bool
	for _, n := range nodes {
		if n.ID == "node-1" && n.Healthy == false {
			updated = true
			break
		}
	}
	if !updated {
		t.Errorf("expected node-1 to be updated to unhealthy")
	}
}

func TestNodeManager_RegisterDuplicate(t *testing.T) {
	ds := NewFakeDatastore()
	manager := NewManager(ds)

	node1 := models.Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node1); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Try to register the same node again with the same Healthy status.
	err := manager.Register(node1)
	if err == nil {
		t.Errorf("expected error when registering duplicate node, got nil")
	}
}
