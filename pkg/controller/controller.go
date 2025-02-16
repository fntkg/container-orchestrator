package controller

import (
	"log"
	"time"

	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
)

// Manager monitors and reconciles the state of the cluster.
type Manager struct {
	scheduler   scheduler.Scheduler
	tasks       []scheduler.Task
	nodeManager node.NodeManager
}

// NewControllerManager creates a new instance of Manager.
func NewControllerManager(scheduler scheduler.Scheduler, tasks []scheduler.Task, nm node.NodeManager) *Manager {
	return &Manager{
		scheduler:   scheduler,
		tasks:       tasks,
		nodeManager: nm,
	}
}

// Run starts the reconciliation loop.
// It listens for a stop signal via the stopCh channel.
func (cm *Manager) Run(stopCh <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.reconcile()
		case <-stopCh:
			log.Println("Controller Manager stopped")
			return
		}
	}
}

// reconcile performs a single reconciliation iteration.
func (cm *Manager) reconcile() {
	log.Println("Controller Manager: Reconciling state...")

	// Retrieve the list of nodes from the Node Manager.
	nodes := cm.nodeManager.GetNodes()
	// filter out unhealthy nodes.
	//healthyNodes := make([]node.Node, 0)
	var healthyNodes []node.Node
	for _, n := range nodes {
		if n.Healthy {
			healthyNodes = append(healthyNodes, n)
		}
	}

	for _, task := range cm.tasks {
		assignedNode, err := cm.scheduler.Schedule(task, healthyNodes)
		if err != nil {
			log.Printf("Error scheduling task %s: %v", task.ID, err)
		} else {
			log.Printf("Task %s assigned to Node %s", task.ID, assignedNode.ID)
		}
	}
}
