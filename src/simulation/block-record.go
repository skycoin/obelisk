package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type BlockRecord struct {
	hash     cipher.SHA256  // Hash of the Block
	seqNo    int            // SeqNo of the block
	parent   *BlockRecord   // Pointer to the parent of the block record
	children []*BlockRecord // List of children of the block record
}

func NewRandomBlockRecord() *BlockRecord {
	blockRecord := &BlockRecord{}
	blockRecord.hash = GetRandomSHA256()
	blockRecord.seqNo = 0
	blockRecord.parent = nil
	blockRecord.children = []*BlockRecord{}
	return blockRecord
}

func NewRandomChildBlockRecord(parent *BlockRecord) *BlockRecord {

	blockRecord := &BlockRecord{}
	blockRecord.hash = GetRandomSHA256()
	blockRecord.seqNo = 0
	blockRecord.parent = parent
	blockRecord.children = []*BlockRecord{}
	return blockRecord
}
