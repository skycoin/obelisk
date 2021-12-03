package main

// import (
// 	"github.com/skycoin/skycoin/src/cipher"
// )

type NodeBlockMeta struct {
	blockRecord *BlockRecord // The corresponding BlockRecord from the tree
	seqNo       int          // Here we maintain the highest seqNo among the nodes considered while syncing the states
	ticks       int          // capture the ticks from the global simulation at the time when the NodeBlockMeta was synced
	weight      float64      // weight of NodeBlockMeta. This weight will be sum of the weight of the children of the blockRecord
}

func (n *NodeBlockMeta) VerifyNodeBlockMeta() {
}

func NewNodeBlockMeta(blockRecord *BlockRecord) *NodeBlockMeta {
	return &NodeBlockMeta{blockRecord: blockRecord, seqNo: 0, ticks: 0, weight: 0}
}
