package main

import (
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher"
)

type Simulation struct {
	Seed          int64
	Iterations    int
	Ticks         int
	Nodes         map[cipher.PubKey]*Node
	RootBlockTree *BlockRecordTree
}

var simulation *Simulation

func GetSimulation() *Simulation {
	if simulation == nil {
		simulation = &Simulation{Nodes: map[cipher.PubKey]*Node{}}
	}
	return simulation
}

func (sim *Simulation) InitSimulation(totalRootBlockTreeNodes int, totalRootBlockTreeChildrenPerNode int, numberOfNodes int,
	numberOfSubscribers int, seed int64, iterations int) error {
	if rootBlockTree, err := NewRandomBlockRecordTree(totalRootBlockTreeNodes, totalRootBlockTreeChildrenPerNode); err != nil {
		return err
	} else {
		sim.RootBlockTree = rootBlockTree
	}

	sim.Nodes = map[cipher.PubKey]*Node{}
	nodes := []*Node{}
	var i int
	for i < numberOfNodes {

		node := NewRandomNode()
		sim.Nodes[node.pubKey] = node
		nodes = append(nodes, node)
		i++
	}

	for _, node := range sim.Nodes {
		node.InitializeNode(sim.RootBlockTree, nodes, numberOfSubscribers, seed)
	}

	sim.Seed = seed
	sim.Ticks = 0
	sim.Iterations = iterations

	return nil
}

func (sim *Simulation) AdvanceTicks() {
	sim.Ticks++
}

func (sim *Simulation) RunSimulation() {

	rand.Seed(sim.Seed)

	nodeArray := []*Node{}
	for _, node := range sim.Nodes {
		nodeArray = append(nodeArray, node)
	}

	var it int = 0
	for it < sim.Iterations {
		nodeArray[rand.Intn(len(sim.Nodes))].UpdateNodeState()
		it++
	}

	sim.Ticks++
}

func (sim *Simulation) PrintTotalState() {
	// Print Total State logic
}
