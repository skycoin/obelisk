package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type BlockRecord struct {
	hash     cipher.SHA256  
	seqNo    int            
	parent   *BlockRecord   
	children []*BlockRecord 
}

// NewRandomBlockRecord: Create a new Block Record with a random hash
func NewRandomBlockRecord() *BlockRecord {
	blockRecord := &BlockRecord{}
	blockRecord.hash = GetRandomSHA256()
	blockRecord.seqNo = 0
	blockRecord.parent = nil
	blockRecord.children = []*BlockRecord{}
	return blockRecord
}
