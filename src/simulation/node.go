package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/skycoin/skycoin/src/cipher"
)

const CONSENSUS_APPROACH_FACTOR = 0.1
const WEIGHT_INIT_FACTOR = 0.01
const NODE_INIT_RANDOM_SUBSCRIPTION_POLICY = "NODE_INIT_RANDOM_SUBSCRIPTION_POLICY"
const NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY = "NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY"

type Node struct {
	id            int                                   
	pubKey        cipher.PubKey                         
	seqNo         int                                   
	subscriptions []*Node                               
	state         map[cipher.SHA256]*NodeStateBlockMeta 
	nodeMessagesReceived []*NodeMessage
	nodeMessageQueue  *NodeMessageQueue
}

// NewRandomNode: Create a new random Node with a random public key
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

// InitializeNode: Initialize a node based on a provided subscriber selection policy
func (n *Node) InitializeNode(sim *Simulation, numberOfSubscribers int) {

	switch sim.NodeInitPolicy {
		case NODE_INIT_RANDOM_SUBSCRIPTION_POLICY:
			n.initializeRandomNodeSubcribers(sim.Nodes, numberOfSubscribers);
		case NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY:
			n.initializeNodeSubscribersViaAnnualRing(sim.NodeGridMap.gridMap[n.pubKey], sim.Nodes, numberOfSubscribers);
	}
	n.InitializeNodeState(sim.RootBlockTree);
}

// Initializes a node's subscribers by randomly choosing subscribers from the overall nodes available
func (n *Node) initializeRandomNodeSubcribers(nodes []*Node, numberOfSubscribers int) {

	reuseCheckMap := map[int]bool{}

	for len(n.subscriptions) < numberOfSubscribers {
		subscriberIndex := rand.Intn(len(nodes))

		if _, ok := reuseCheckMap[subscriberIndex]; !ok && n != nodes[subscriberIndex] {
			reuseCheckMap[subscriberIndex] = true
			n.subscriptions = append(n.subscriptions, nodes[subscriberIndex])
		}
	}
}

// Initializes a node's subscribers via Annual Ring Algorithm using Node Grids
func (n *Node) initializeNodeSubscribersViaAnnualRing(nodeGrid *NodeGrid, nodes []*Node, numberOfSubscribers int) {
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
				possibleSubscriber := nodeGrid.getValue(i,j);
				if(len(subscribers) < numberOfSubscribers && possibleSubscriber != nil && !checkIfExists(possibleSubscriber, subscribers)) {
					normDistance := math.Sqrt(float64((i * i) + (j * j)));
					if((normDistance >= r1 && normDistance <= r2)) {
						subscribers = append(subscribers, possibleSubscriber);
					}
				}
			}
		}
	}

	n.subscriptions = subscribers;
}

func checkIfExists(possibleSubscriber *Node, subscribers []*Node) bool {

	if(possibleSubscriber == nil || len(subscribers) <= 0) {
		return false;
	}

	for _, subscriber := range subscribers {
		if(possibleSubscriber == subscriber) {
			return true;
		}
	}

	return false;
}

// InitializeNodeState: Initialize a node's state by creating a block record tree and assigning weights 
func (n *Node) InitializeNodeState(brt *BlockRecordTree) {
	blockRecords := brt.GetAllBlockRecords()

	for _, blockRecord := range blockRecords {
		n.state[blockRecord.hash] = NewNodeStateBlockMeta(blockRecord)
	}

	n.setWeight(1.0, brt.Root);

}

func getRandomSignMultiplier() float64 {
	multiplier := 1.0;
	if (rand.Intn(2) == 0) {
		multiplier = -1.0
	}
	return multiplier;
}

// ValidateNodeState: Validate a node's state by checking that each parent block's weight is equal to the sum of the weights of the children 
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

// UpdateNodeState: 
// 1) Receive Messages available in the current tick 
// 2) Compute weights of a node's state from it's subcribers
// 3) Send messages to subscribers
func (n *Node) UpdateNodeState() {
	sim := GetSimulation()
	sim.AdvanceTicks()
	n.seqNo++
	
	n.nodeMessagesReceived = append(n.nodeMessagesReceived, n.nodeMessageQueue.Pop(GetSimulation().Ticks)...);

	for _, blockStateMeta := range n.state {
		blockStateMeta.seqNo = n.getMaxSubscribersSeqNo(blockStateMeta.blockRecord.hash)
		n.state[blockStateMeta.blockRecord.hash].weight = n.calculateNewBlockStateMetaWeight(blockStateMeta.blockRecord);
		blockStateMeta.ticks = sim.Ticks
	}
	
	n.adjustWeightsTowardsConsensus(sim.RootBlockTree.Root);

	for _, subscription := range n.subscriptions {
		subscription.sendMessage(n, sim.Ticks+sim.CommunicationsDelayMatrix.matrix[n.pubKey][subscription.pubKey], fmt.Sprintf("Hello from Node %d to Node %d", n.id, subscription.id));
	}
}

// get the max sequence number among all subscribers of the node for a particular block record hash
func (n *Node) getMaxSubscribersSeqNo(hash cipher.SHA256) int {
	var maxSeqNo int = 0

	for _, subscription := range n.subscriptions {
		if _, ok := subscription.state[hash]; ok && subscription.seqNo > maxSeqNo {
			maxSeqNo = subscription.seqNo
		}
	}

	return maxSeqNo
}

// calculate a node's weight for a particular block record by averaging the weights of the same block across all subcribers
func (n *Node) calculateNewBlockStateMetaWeight(blockRecord *BlockRecord) float64 {

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

// PrintNodeDetails: Prints a node's details (id, seqNo, PubKey, Subscriptions, Messages Received, State) 
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
	
	for _, blockRecord := range GetSimulation().RootBlockTree.GetAllBlockRecords() {
		var parentHash cipher.SHA256
		if blockRecord.parent != nil {
			parentHash = blockRecord.parent.hash
		}
		fmt.Printf("%v | %v | %d | %d | %.2f\n", blockRecord.hash, parentHash, 
		n.state[blockRecord.hash].seqNo, n.state[blockRecord.hash].ticks, n.state[blockRecord.hash].weight)
	}

}

// adjust a node's weight for it to approach a consensus
func (n *Node) adjustWeightsTowardsConsensus(root *BlockRecord) {
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
		n.adjustWeightsTowardsConsensus(child);
	}
}

// set node's weight for a particular block record
func (n *Node) setWeight(newWeight float64, blockRecord *BlockRecord) {

	n.state[blockRecord.hash].weight = newWeight;

	runningWeight := newWeight;
	avgWeight := newWeight / float64(len(blockRecord.children));

	for _, child := range blockRecord.children {
		if(child == blockRecord.children[len(blockRecord.children)-1]) {
			n.setWeight(runningWeight, child);

		} else if (runningWeight > avgWeight) {
			assignWeight := avgWeight + getRandomSignMultiplier() * WEIGHT_INIT_FACTOR
			n.setWeight(assignWeight, child);
			runningWeight -= assignWeight;
		} else {
			n.setWeight(runningWeight, child)
		}
	}
}

// push a message toa node's message queue while settings the desired arrival time
func (n *Node) sendMessage(from *Node, arrivalTimeTick int, message string) {
	n.nodeMessageQueue.Push(&NodeMessage{
		arrivalTimeTick: arrivalTimeTick,
		sentTimeTick: GetSimulation().Ticks,
		from: from,
		to: n,
		message: message,
	});
}
