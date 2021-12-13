package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
)

// The rate at which the nodes should approach the consensus
const CONSENSUS_APPROACH_FACTOR = 0.1
const WEIGHT_INIT_FACTOR = 0.01

type Node struct {
	id            int                                   // id of the node
	pubKey        cipher.PubKey                         // Node's public key
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

func (n *Node) InitializeNode(brt *BlockRecordTree, nodes []*Node, numberOfSubscribers int) {
	n.InitializeRandomNodeSubcribers(nodes, numberOfSubscribers)
	n.InitializeNodeState(brt)
}

func (n *Node) InitializeRandomNodeSubcribers(nodes []*Node, numberOfSubscribers int) {

	reuseCheckMap := map[int]bool{}

	for len(n.subscriptions) < numberOfSubscribers {
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
	}

	n.SetWeight(1.0, brt.Root);

}

func getRandomSignMultiplier() float64 {
	multiplier := 1.0;
	if (rand.Intn(2) == 0) {
		multiplier = -1.0
	}
	return multiplier;
}

func (n *Node) ValidateNodeState() error {
	for _, blockStateMeta := range n.state {
		if len(blockStateMeta.blockRecord.children) > 0 {
			parentWeight := blockStateMeta.weight

			childrenSum := 0.0
			for _, child := range blockStateMeta.blockRecord.children {
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

	for _, blockStateMeta := range n.state {
		blockStateMeta.seqNo = n.GetMaxSubscribersSeqNo(blockStateMeta.blockRecord.hash)
		n.state[blockStateMeta.blockRecord.hash].weight = n.CalculateNewBlockStateMetaWeight(blockStateMeta.blockRecord);
		blockStateMeta.ticks = sim.Ticks
	}

	// Adjust weights towards consensus
	n.AdjustWeightsTowardsConsensus(GetSimulation().RootBlockTree.Root);
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

func (n *Node) CalculateNewBlockStateMetaWeight(blockRecord *BlockRecord) float64 {

	totalWeight := 0.0
	subscriberBlockCount := 0.0

	for _, subscription := range n.subscriptions {
		if _, ok := subscription.state[blockRecord.hash]; ok {
			totalWeight += subscription.state[blockRecord.hash].weight
			subscriberBlockCount++
		}
	}

	var subscriberBlockWeightAvg = 0.0; 

	if subscriberBlockCount > 0 {
		subscriberBlockWeightAvg = totalWeight / subscriberBlockCount
	}

	return subscriberBlockWeightAvg;
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

func (n *Node) AdjustWeightsTowardsConsensus(root *BlockRecord) {
	totalWeight := n.state[root.hash].weight;
	runningWeight := totalWeight

	for _, child := range root.children {
		
		if(child == root.children[len(root.children)-1]) {
			n.state[child.hash].weight = runningWeight;
		} else if(n.state[child.hash].weight > (totalWeight / float64(len(root.children)))) {
			n.state[child.hash].weight += CONSENSUS_APPROACH_FACTOR;
			if(n.state[child.hash].weight > runningWeight) {
				n.state[child.hash].weight = runningWeight;
			}
		} else {
			n.state[child.hash].weight -= CONSENSUS_APPROACH_FACTOR;
			if(n.state[child.hash].weight < 0.0) {
				n.state[child.hash].weight = 0.0;
			}
		}

		runningWeight -= n.state[child.hash].weight;
		n.AdjustWeightsTowardsConsensus(child);
	}
}

func (n *Node) SetWeight(newWeight float64, blockRecord *BlockRecord) {

	n.state[blockRecord.hash].weight = newWeight;

	// Re-assign children weight
	runningWeight := newWeight;
	avgWeight := newWeight / float64(len(blockRecord.children));

	for _, child := range blockRecord.children {
		if(child == blockRecord.children[len(blockRecord.children)-1]) {
			n.SetWeight(runningWeight, child);

		} else if (runningWeight > avgWeight) {
			assignWeight := avgWeight + getRandomSignMultiplier() * WEIGHT_INIT_FACTOR
			n.SetWeight(assignWeight, child);
			runningWeight -= assignWeight;
		} else {
			n.SetWeight(runningWeight, child)
		}
	}
}
