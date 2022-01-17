package main

import "math/rand"

// NodeGrid: A Grid of nodes showing their distances to a particular node
type NodeGrid struct {
	grid [][]*Node
}

// InitializeNodeGrid: Initializes NodeGrid for a given nodes and places the rest of the nodes randomly on the grid
func (ng *NodeGrid) InitializeNodeGrid(initNode *Node, nodes []*Node) {

	ng.grid = [][]*Node{};
	for i := 0; i < len(nodes); i++ {
		nodeList := []*Node{};
		for j := 0; j < len(nodes); j++ { 
			nodeList = append(nodeList, nil);
		}
		ng.grid = append(ng.grid, nodeList)
	};

	nodesToAssign := []*Node{};
	for _, node := range nodes {
		if(node != initNode) {
			nodesToAssign = append(nodesToAssign, node)
		}
	}

	for len(nodesToAssign) > 0 {
		i := rand.Intn(len(nodes));
		j := rand.Intn(len(nodes));

		if(ng.grid[i][j] == nil) {
			ng.grid[i][j] = nodesToAssign[len(nodesToAssign)-1];
			nodesToAssign = nodesToAssign[:len(nodesToAssign)-1];
		}
	}
}

func (ng *NodeGrid) getLength() int {
	return len(ng.grid);
}

func (ng *NodeGrid) getValue(i, j int) *Node {
	return ng.grid[i][j];
}
