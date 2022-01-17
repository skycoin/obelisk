package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

// NodeGridMap: Keeps a mapping for Each node's PubKey to it's NodeGrid
type NodeGridMap struct {
	gridMap map[cipher.PubKey]*NodeGrid
}

func (ngm *NodeGridMap) InitializeNodeGridMap(nodes []*Node) {
	ngm.gridMap = map[cipher.PubKey]*NodeGrid{}
	for _, node := range nodes {
		nodeGrid := &NodeGrid{};
		nodeGrid.InitializeNodeGrid(node, nodes)
		ngm.gridMap[node.pubKey] = nodeGrid; 
	}
}
