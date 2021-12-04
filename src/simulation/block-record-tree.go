package main

import (
	"fmt"
)

type BlockRecordTree struct {
	Root *BlockRecord
}

func (brt *BlockRecordTree) GetAllBlockRecords() []*BlockRecord {
	blockRecordArray := []*BlockRecord{}

	// queue is needed for breadth first tree traversal
	queue := []*BlockRecord{}
	queue = append(queue, brt.Root)

	for len(queue) > 0 {
		blockRecordArray = append(blockRecordArray, queue[0])
		for _, child := range queue[0].children {
			queue = append(queue, child)
		}

		queue = queue[1:]
	}

	return blockRecordArray
}

func NewRandomBlockRecordTree(totalBlocks int, childrenPerNode int) (*BlockRecordTree, error) {
	if totalBlocks < 1 {
		return nil, fmt.Errorf("totalBlocks must be greater than 0")
	}

	if childrenPerNode < 1 {
		return nil, fmt.Errorf("childrenPerNode must be greater than 0")
	}

	blockRecordTree := &BlockRecordTree{}
	blockRecordTree.Root = NewRandomBlockRecord()

	// queue is needed for breadth first tree traversal
	queue := []*BlockRecord{}
	queue = append(queue, blockRecordTree.Root)
	totalBlocks--

	var n int = 0
	for n < totalBlocks {
		currentRoot := queue[0]

		br := NewRandomBlockRecord()

		if len(currentRoot.children) < childrenPerNode {
			br.parent = currentRoot
			br.seqNo = br.parent.seqNo + 1
			currentRoot.children = append(currentRoot.children, br)
		} else {
			queue = queue[1:]
			currentRoot = queue[0]
			br.parent = currentRoot
			br.seqNo = br.parent.seqNo + 1
			currentRoot.children = append(currentRoot.children, br)
		}

		queue = append(queue, br)

		n++
	}

	return blockRecordTree, nil
}
