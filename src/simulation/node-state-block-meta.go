package main

type NodeStateBlockMeta struct {
	blockRecord *BlockRecord // The corresponding BlockRecord from the tree
	seqNo       int          // Here we maintain the highest seqNo among the nodes considered while syncing the states
	ticks       int          // capture the ticks from the global simulation at the time when the NodeStateBlockMeta was synced
	weight      float64      // weight of NodeStateBlockMeta. This weight will be sum of the weight of the children of the blockRecord
}

func (n *NodeStateBlockMeta) VerifyNodeStateBlockMeta() {
}

func NewNodeStateBlockMeta(blockRecord *BlockRecord) *NodeStateBlockMeta {
	return &NodeStateBlockMeta{blockRecord: blockRecord, seqNo: 0, ticks: 0, weight: 0}
}
