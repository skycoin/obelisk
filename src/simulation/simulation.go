package main

import (
	"fmt"
	"math/rand"
)

type Simulation struct {
	VerboseMode   bool
	Iterations    int
	Ticks         int
	Nodes         []*Node
	RootBlockTree *BlockRecordTree
}

var simulation *Simulation

func GetSimulation() *Simulation {
	if simulation == nil {
		simulation = &Simulation{}
	}
	return simulation
}

func (sim *Simulation) InitSimulation(totalRootBlockTreeNodes int, totalRootBlockTreeChildrenPerNode int, numberOfNodes int,
	numberOfSubscribers int, iterations int, verboseMode bool) error {
	if rootBlockTree, err := NewRandomBlockRecordTree(totalRootBlockTreeNodes, totalRootBlockTreeChildrenPerNode); err != nil {
		return err
	} else {
		sim.RootBlockTree = rootBlockTree
	}

	sim.Nodes = []*Node{}
	for i := 0; i < numberOfNodes; i++ {
		node := NewRandomNode(i + 1)
		sim.Nodes = append(sim.Nodes, node)
	}

	for _, node := range sim.Nodes {
		node.InitializeNode(sim.RootBlockTree, sim.Nodes, numberOfSubscribers)
	}

	sim.Ticks = 0
	sim.Iterations = iterations
	sim.VerboseMode = verboseMode

	return nil
}

func (sim *Simulation) AdvanceTicks() {
	sim.Ticks++
}

func (sim *Simulation) RunSimulation() error {

	// Printing before simulation
	fmt.Printf("\n\n\n#begin Simulation Initial State:\n")
	sim.PrintAllNodes()
	fmt.Printf("\n#end Simulation Initial State\n")

	for it := 0; it < sim.Iterations; it++ {

		node := sim.Nodes[rand.Intn(len(sim.Nodes))]

		if sim.VerboseMode {
			fmt.Printf("\n\n\nIteration No. %d", (it + 1))
			fmt.Printf("\n\nBefore Update:\n")
			node.PrintNodeDetails()
		}

		node.UpdateNodeState()

		if err := node.ValidateNodeState(); err != nil {
			return fmt.Errorf("Node State Validation failed for node-id:%d err:%s", node.id, err.Error())
		}

		if sim.VerboseMode {
			fmt.Printf("\n\nAfter Update:\n")
			node.PrintNodeDetails()
		}
	}

	sim.Ticks++

	// Printing after simulation
	fmt.Printf("\n\n\n#begin Simulation Final State:\n")
	sim.PrintAllNodes()
	fmt.Printf("\n#end Simulation Final State\n")

	return nil
}

func (sim *Simulation) PrintAllNodes() {

	for _, node := range sim.Nodes {
		fmt.Printf("\n")
		node.PrintNodeDetails()
	}
}
