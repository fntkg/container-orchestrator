package node

import (
	"errors"
	"github.com/fntkg/container-orchestrator/pkg/models"

	"github.com/fntkg/container-orchestrator/pkg/datastore"
)

// NodeManager is an interface that defines the behavior of a node manager.
type NodeManager interface {
	Register(n models.Node) error
	GetNodes() []models.Node
	UpdateHealth(nodeID string, healthy bool) error
}

// DefaultNodeManager manages the nodes in the cluster.
type DefaultNodeManager struct {
	ds datastore.Datastore
}

// NewManager // NewManager creates a new instance of DefaultNodeManager with the given datastore.
func NewManager(ds datastore.Datastore) *DefaultNodeManager {
	return &DefaultNodeManager{
		ds: ds,
	}
}

// Register adds a new node to the manager.
func (m *DefaultNodeManager) Register(n models.Node) error {
	return m.ds.SaveNode(n)
}

// GetNodes returns a slice of all registered nodes.
func (m *DefaultNodeManager) GetNodes() []models.Node {
	nodes, err := m.ds.GetNodes()
	if err != nil {
		// Optionally, you can log the error here
		return []models.Node{}
	}
	return nodes
}

// UpdateHealth updates the health status of a node.
func (m *DefaultNodeManager) UpdateHealth(nodeID string, healthy bool) error {
	// Retrieve all nodes from the datastore.
	nodes, err := m.ds.GetNodes()
	if err != nil {
		return err
	}

	// Find the node with the given nodeID.
	var found bool
	var updatedNode models.Node
	for _, n := range nodes {
		if n.ID == nodeID {
			updatedNode = n
			updatedNode.Healthy = healthy
			found = true
			break
		}
	}
	if !found {
		return errors.New("node not found")
	}

	// Save the updated node back to the datastore.
	return m.ds.SaveNode(updatedNode)
}
