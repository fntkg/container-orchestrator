// File: pkg/controller/controller.go
package controller

import (
	"log"
	"time"

	"github.com/fntkg/container-orchestrator/pkg/models"
	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
	"github.com/fntkg/container-orchestrator/pkg/taskmanager"
)

// ControllerManager monitors and reconciles the cluster's state.
type ControllerManager struct {
	scheduler   scheduler.Scheduler
	taskManager taskmanager.TaskManager
	nodeManager node.NodeManager
}

// NewControllerManager creates a new ControllerManager instance.
func NewControllerManager(sched scheduler.Scheduler, tm taskmanager.TaskManager, nm node.NodeManager) *ControllerManager {
	return &ControllerManager{
		scheduler:   sched,
		taskManager: tm,
		nodeManager: nm,
	}
}

// Run starts the reconciliation loop.
func (cm *ControllerManager) Run(stopCh <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.reconcile()
		case <-stopCh:
			log.Println("Controller DefaultNodeManager stopped")
			return
		}
	}
}

// reconcile performs a single reconciliation iteration.
func (cm *ControllerManager) reconcile() {
	log.Println("Controller DefaultNodeManager: Reconciling state...")

	// Retrieve tasks from Task DefaultNodeManager.
	tasks, err := cm.taskManager.GetTasks()
	if err != nil {
		log.Printf("Error retrieving tasks: %v", err)
		return
	}

	// Retrieve nodes from Node DefaultNodeManager.
	nodes := cm.nodeManager.GetNodes()
	// Filter only healthy nodes.
	healthyNodes := make([]models.Node, 0)
	for _, n := range nodes {
		if n.Healthy {
			healthyNodes = append(healthyNodes, n)
		}
	}

	for _, task := range tasks {
		assignedNode, err := cm.scheduler.Schedule(task, healthyNodes)
		if err != nil {
			log.Printf("Error scheduling task %s: %v", task.ID, err)
		} else {
			log.Printf("Task %s assigned to Node %s", task.ID, assignedNode.ID)
			// Optionally, update task status via Task DefaultNodeManager.
		}
	}
}
