package main

import (
	"fmt"
	"math/rand"
)

const DEFAULT_NODE_INIT_POLICY = NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY;

type Simulation struct {
	VerboseMode   	bool
	RootBlockTree 	*BlockRecordTree
	Iterations    	int
	Ticks         	int
	NodeInitPolicy  string
	NodeGridMap   	*NodeGridMap
	Nodes         	[]*Node
	CommunicationsDelayMatrix *CommunicationsDelayMatrix
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

	sim.CommunicationsDelayMatrix = &CommunicationsDelayMatrix{}
	sim.CommunicationsDelayMatrix.InitializeCommunicationsDelayMatrix(sim.Nodes);

	sim.NodeGridMap = &NodeGridMap{};
	sim.NodeGridMap.InitializeNodeGridMap(sim.Nodes);
	sim.NodeInitPolicy = DEFAULT_NODE_INIT_POLICY;

	for _, node := range sim.Nodes {
		node.InitializeNode(sim, numberOfSubscribers)
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

	convergenceAchievedIteration := -1;

	for it := 0; it < sim.Iterations; it++ {

		node := sim.Nodes[rand.Intn(len(sim.Nodes))]

		if sim.VerboseMode {
			fmt.Printf("\n\n\nIteration No. %d", (it + 1))
			fmt.Printf("\n\nBefore Update:\n")
			node.PrintNodeDetails()
		}

		node.UpdateNodeState()

		if err := node.ValidateNodeState(); err != nil {
			return fmt.Errorf("Node State Validation failed for node-id:%d Iteration:%d err:%s", node.id, it, err.Error());
		}

		if sim.VerboseMode {
			fmt.Printf("\n\nAfter Update:\n");
			node.PrintNodeDetails()
		}

		if sim.checkConvergence() {
			convergenceAchievedIteration = it;
			break;
		}
	}

	sim.Ticks++

	// Printing after simulation
	fmt.Printf("\n\n\n#begin Simulation Final State:\n")
	sim.PrintAllNodes()
	fmt.Printf("\n#end Simulation Final State\n")

	if(convergenceAchievedIteration != -1) {
		fmt.Printf("\nIteration %d: Convergence Achieved!!!", convergenceAchievedIteration)
	}

	return nil
}

func (sim *Simulation) PrintAllNodes() {

	for _, node := range sim.Nodes {
		fmt.Printf("\n")
		node.PrintNodeDetails()
	}
}

func (sim *Simulation) checkConvergence() bool {

	allBlocks := sim.RootBlockTree.GetAllBlockRecords();

	for _, block := range allBlocks {

		firstWeight := sim.Nodes[0].state[block.hash].weight;		
		
		// Verify that for all nodes the weight of the block is same
		for _, node := range sim.Nodes {
			if(node.state[block.hash].weight != firstWeight) {
				return false;
			}
		}

		// verify that the weight value is either 0.0 or 1.0
		if  (firstWeight != 0.0 && firstWeight != 1.0) {
			return false;
		}
	} 

	return true;
}
