package node

import "testing"

func TestNodeManager_RegisterAndGetNodes(t *testing.T) {
	manager := NewManager()

	// Register a couple of nodes.
	nodesToRegister := []Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: true},
	}

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
	manager := NewManager()

	// Register a node.
	node := Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Update health status.
	if err := manager.UpdateHealth("node-1", false); err != nil {
		t.Fatalf("failed to update node health: %v", err)
	}

	// Verify the update.
	nodes := manager.GetNodes()
	for _, n := range nodes {
		if n.ID == "node-1" && n.Healthy != false {
			t.Errorf("expected node-1 to be unhealthy, but got healthy")
		}
	}
}

func TestNodeManager_RegisterDuplicate(t *testing.T) {
	manager := NewManager()

	node := Node{ID: "node-1", Healthy: true}
	if err := manager.Register(node); err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Try to register the same node again.
	err := manager.Register(node)
	if err == nil {
		t.Errorf("expected error when registering duplicate node, got nil")
	}
}
