package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher"
)

type Node struct {
	pubKey        cipher.PubKey
	seqNo         int                              // Node's sequence number tracking the number of updates done on the node
	subscriptions []*Node                          // List of Nodes subscribed by the current Node
	state         map[cipher.SHA256]*NodeBlockMeta // A mapping from BlockRecord Hash to current Node's separate copy of NodeBlockMeta
}

func NewRandomNode() *Node {
	node := &Node{}
	node.pubKey = GetRandomPubKey()
	node.seqNo = 0
	node.subscriptions = []*Node{}
	node.state = map[cipher.SHA256]*NodeBlockMeta{}

	return node
}

func (n *Node) InitializeNode(brt *BlockRecordTree, nodes []*Node, numberOfSubcribers int, seed int64) {
	n.InitializeRandomNodeSubcribers(nodes, numberOfSubcribers, seed)
	n.InitializeNodeState(brt)
}

func (n *Node) InitializeRandomNodeSubcribers(nodes []*Node, numberOfSubcribers int, seed int64) {

	rand.Seed(seed)
	reuseCheckMap := map[int]bool{}

	for int(len(n.subscriptions)) < numberOfSubcribers {
		subscriberIndex := rand.Intn(int(len(nodes)))

		if _, ok := reuseCheckMap[subscriberIndex]; !ok && n != nodes[subscriberIndex] {
			n.subscriptions = append(n.subscriptions, nodes[subscriberIndex])
		}
	}
}

func (n *Node) InitializeNodeState(brt *BlockRecordTree) {
	blockRecords := brt.GetAllBlockRecords()

	for _, blockRecord := range blockRecords {
		n.state[blockRecord.hash] = NewNodeBlockMeta(blockRecord)

		if brt.Root == blockRecord {
			n.state[blockRecord.hash].weight = 1
		} else {
			n.state[blockRecord.hash].weight = n.state[blockRecord.parent.hash].weight / float64(len(blockRecord.parent.children))
		}
	}
}

func (n *Node) ValidateNodeState() error {
	for _, nodeStateMeta := range n.state {
		if len(nodeStateMeta.blockRecord.children) > 0 {
			parentWeight := nodeStateMeta.weight

			childrenSum := 0.0
			for _, child := range nodeStateMeta.blockRecord.children {
				childrenSum += n.state[child.hash].weight
			}

			if math.Abs(parentWeight-childrenSum) > 0.0000001 {
				return fmt.Errorf("Node %v has weight %f while children sum is %f", n.pubKey, parentWeight, childrenSum)
			}
		}
	}
	return nil
}

func (n *Node) UpdateNodeState() {
	sim := GetSimulation()
	sim.AdvanceTicks()
	n.seqNo++

	for _, nodeStateMeta := range n.state {
		nodeStateMeta.seqNo = n.GetMaxSubscribersSeqNo(nodeStateMeta.blockRecord.hash)
		nodeStateMeta.weight = n.GetAvgSubscribersWeight(nodeStateMeta.blockRecord.hash)
		nodeStateMeta.ticks = sim.Ticks
	}
}

func (n *Node) GetMaxSubscribersSeqNo(hash cipher.SHA256) int {
	var maxSeqNo int = 0

	for _, subscription := range n.subscriptions {
		if _, ok := subscription.state[hash]; ok && subscription.seqNo > maxSeqNo {
			maxSeqNo = subscription.seqNo
		}
	}

	return maxSeqNo
}

func (n *Node) GetAvgSubscribersWeight(hash cipher.SHA256) float64 {
	totalWeight := 0.0
	totalCount := 0.0

	for _, subscription := range n.subscriptions {
		if _, ok := subscription.state[hash]; ok {
			totalWeight += subscription.state[hash].weight
			totalCount++
		}
	}

	if totalCount > 0 {
		return totalWeight / totalCount
	} else {
		return 0.0
	}
}
