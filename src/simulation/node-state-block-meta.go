package main

type NodeStateBlockMeta struct {
	blockRecord *BlockRecord
	seqNo       int
	ticks       int
	weight      float64
}

// NewNodeStateBlockMeta: Creates a NodeStateBlockMeta for a given block record
func NewNodeStateBlockMeta(blockRecord *BlockRecord) *NodeStateBlockMeta {
	return &NodeStateBlockMeta{blockRecord: blockRecord, seqNo: 0, ticks: 0, weight: 0}
}
