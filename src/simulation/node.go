package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
)

type Node struct {
	id            int
	pubKey        cipher.PubKey
	seqNo         int                                   // Node's sequence number tracking the number of updates done on the node
	subscriptions []*Node                               // List of Nodes subscribed by the current Node
	state         map[cipher.SHA256]*NodeStateBlockMeta // A mapping from BlockRecord Hash to current Node's separate copy of NodeStateBlockMeta
}

func NewRandomNode(id int) *Node {
	node := &Node{id: id}
	node.pubKey = GetRandomPubKey()
	node.seqNo = 0
	node.subscriptions = []*Node{}
	node.state = map[cipher.SHA256]*NodeStateBlockMeta{}

	return node
}

func (n *Node) InitializeNode(brt *BlockRecordTree, nodes []*Node, numberOfSubcribers int) {
	n.InitializeRandomNodeSubcribers(nodes, numberOfSubcribers)
	n.InitializeNodeState(brt)
}

func (n *Node) InitializeRandomNodeSubcribers(nodes []*Node, numberOfSubcribers int) {

	reuseCheckMap := map[int]bool{}

	for len(n.subscriptions) < numberOfSubcribers {
		subscriberIndex := rand.Intn(len(nodes))

		if _, ok := reuseCheckMap[subscriberIndex]; !ok && n != nodes[subscriberIndex] {
			reuseCheckMap[subscriberIndex] = true
			n.subscriptions = append(n.subscriptions, nodes[subscriberIndex])
		}
	}
}

func (n *Node) InitializeNodeState(brt *BlockRecordTree) {
	blockRecords := brt.GetAllBlockRecords()

	for _, blockRecord := range blockRecords {
		n.state[blockRecord.hash] = NewNodeStateBlockMeta(blockRecord)

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

func (n *Node) PrintNodeDetails() {
	subscriptionIds := []int{}

	for _, subscription := range n.subscriptions {
		subscriptionIds = append(subscriptionIds, subscription.id)
	}

	fmt.Printf("Node (id=%d seqNo=%d) Details:\n", n.id, n.seqNo)
	fmt.Printf("PubKey:%v\n", n.pubKey)
	fmt.Printf("Subscriptions:%v\n", subscriptionIds)
	fmt.Println("State [Format: blockHash | parentHash | seqNo | ticks | weight]:")

	// Getting sorted states in the descending order of the weights
	sortedStates := []*NodeStateBlockMeta{}
	for _, blockStateMeta := range n.state {
		sortedStates = append(sortedStates, blockStateMeta)
	}
	sort.Slice(sortedStates, func(i, j int) bool { return sortedStates[i].weight > sortedStates[j].weight })

	for _, blockStateMeta := range sortedStates {
		var parentHash cipher.SHA256
		if blockStateMeta.blockRecord.parent != nil {
			parentHash = blockStateMeta.blockRecord.parent.hash
		}
		fmt.Printf("%v | %v | %d | %d | %.2f\n", blockStateMeta.blockRecord.hash, parentHash, blockStateMeta.seqNo, blockStateMeta.ticks, blockStateMeta.weight)
	}
}
