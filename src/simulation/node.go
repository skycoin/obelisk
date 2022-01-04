package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher"
)

// The rate at which the nodes should approach the consensus
const CONSENSUS_APPROACH_FACTOR = 0.1
const WEIGHT_INIT_FACTOR = 0.01
const NODE_INIT_RANDOM_SUBSCRIPTION_POLICY = "NODE_INIT_RANDOM_SUBSCRIPTION_POLICY"
const NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY = "NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY"

type Node struct {
	id            int                                   // id of the node
	pubKey        cipher.PubKey                         // Node's public key
	seqNo         int                                   // Node's sequence number tracking the number of updates done on the node
	subscriptions []*Node                               // List of Nodes subscribed by the current Node
	state         map[cipher.SHA256]*NodeStateBlockMeta // A mapping from BlockRecord Hash to current Node's separate copy of NodeStateBlockMeta
	nodeMessagesReceived []*NodeMessage
	nodeMessageQueue  *NodeMessageQueue
}

func NewRandomNode(id int) *Node {
	node := &Node{id: id}
	node.pubKey = GetRandomPubKey()
	node.seqNo = 0
	node.subscriptions = []*Node{}
	node.state = map[cipher.SHA256]*NodeStateBlockMeta{}
	node.nodeMessagesReceived = []*NodeMessage{}
	node.nodeMessageQueue = &NodeMessageQueue{}
	return node
}

func (n *Node) InitializeNode(sim *Simulation, numberOfSubscribers int) {

	switch sim.NodeInitPolicy {
		case NODE_INIT_RANDOM_SUBSCRIPTION_POLICY:
			n.InitializeRandomNodeSubcribers(sim.Nodes, numberOfSubscribers)
			break;
		case NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY:
			n.InitializeNodeSubscribersViaAnnualRing(sim.NodeGridMap.gridMap[n.pubKey], sim.Nodes, numberOfSubscribers)
			break;
	}
	n.InitializeNodeState(sim.RootBlockTree);
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

func (n *Node) InitializeNodeSubscribersViaAnnualRing(nodeGrid *NodeGrid, nodes []*Node, numberOfSubscribers int) {
	gridLength := nodeGrid.getLength()
	radius := float64(gridLength);
	ringWidthRatio := 0.0;
	subscribers := []*Node{};

	for len(subscribers) < numberOfSubscribers {
		ringWidthRatio += 0.1;
		d := radius * ringWidthRatio;
		r1 := radius - (d / 2);  
		r2 := radius + (d / 2);  
		for i := 0; i < gridLength; i++ {
			for j := 0; j < gridLength; j++ { 
				if(len(subscribers) < numberOfSubscribers && nodeGrid.getValue(i,j) != nil) {
					normDistance := math.Sqrt(float64((i * i) + (j * j)));
					if((normDistance >= r1 && normDistance <= r2)) {
						subscribers = append(subscribers, nodeGrid.getValue(i,j));
					}
				}
			}
		}
	}

	n.subscriptions = subscribers;
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

	// Receiving Messages
	n.nodeMessagesReceived = append(n.nodeMessagesReceived, n.nodeMessageQueue.Pop(GetSimulation().Ticks)...);

	for _, blockStateMeta := range n.state {
		blockStateMeta.seqNo = n.GetMaxSubscribersSeqNo(blockStateMeta.blockRecord.hash)
		n.state[blockStateMeta.blockRecord.hash].weight = n.CalculateNewBlockStateMetaWeight(blockStateMeta.blockRecord);
		blockStateMeta.ticks = sim.Ticks
	}

	// Adjust weights towards consensus
	n.AdjustWeightsTowardsConsensus(sim.RootBlockTree.Root);

	// Sending Messages to subscribers
	for _, subscription := range n.subscriptions {
		subscription.SendMessage(n, sim.Ticks+sim.CommunicationsDelayMatrix.matrix[n.pubKey][subscription.pubKey], fmt.Sprintf("Hello from Node %d to Node %d", n.id, subscription.id));
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
	fmt.Printf("PubKey:%s\n", n.pubKey.Hex())
	fmt.Printf("Subscriptions:%v\n", subscriptionIds)
	if(len(n.nodeMessagesReceived) > 0) {
		fmt.Println("Messages Received [Format: from(pubKey) | to(pubKey) | sentTick | arrivedTick | message]:")
		for _, nodeMessage := range n.nodeMessagesReceived {
			fmt.Printf("%s | %s | %d | %d | %s\n", string(nodeMessage.from.pubKey.Hex()), string(nodeMessage.to.pubKey.Hex()), nodeMessage.sentTimeTick, nodeMessage.arrivalTimeTick, nodeMessage.message)  
		}
	}

	fmt.Println("State [Format: blockHash | parentHash | seqNo | ticks | weight]:")

	// Note State Blocks will be printed in the order Breadth first search tree traversal
	for _, blockRecord := range GetSimulation().RootBlockTree.GetAllBlockRecords() {
		var parentHash cipher.SHA256
		if blockRecord.parent != nil {
			parentHash = blockRecord.parent.hash
		}
		fmt.Printf("%v | %v | %d | %d | %.2f\n", blockRecord.hash, parentHash, 
		n.state[blockRecord.hash].seqNo, n.state[blockRecord.hash].ticks, n.state[blockRecord.hash].weight)
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

func (n *Node) SendMessage(from *Node, arrivalTimeTick int, message string) {
	n.nodeMessageQueue.Push(&NodeMessage{
		arrivalTimeTick: arrivalTimeTick,
		sentTimeTick: GetSimulation().Ticks,
		from: from,
		to: n,
		message: message,
	});
}
