package main

import (
	"sort"
)

// NodeMessageQueue: Keeps track of the messages to be received by a node sorted in the order of arrival time tick
type NodeMessageQueue struct {
	messages []*NodeMessage
}

// Push: Push a message to nodeMessage while ensuring the arrival time order
func (nmq *NodeMessageQueue) Push(nodeMessage *NodeMessage) {
	nmq.messages = append(nmq.messages, nodeMessage);
	sort.SliceStable(nmq.messages, func(i, j int) bool {
		return nmq.messages[i].arrivalTimeTick < nmq.messages[j].arrivalTimeTick
	})
}

// Pop: Pop a message from the start of the queue
func (nmq *NodeMessageQueue) Pop(currentTick int) []*NodeMessage {
	receivedMessages := []*NodeMessage{};

	for len(nmq.messages) > 0 && nmq.messages[0].arrivalTimeTick <= currentTick {
		receivedMessages = append(receivedMessages, nmq.messages[0]);
		if(len(nmq.messages) > 1) {
			nmq.messages = nmq.messages[1:];
		} else {
			nmq.messages = []*NodeMessage{};
		}
	}  

	return receivedMessages;
}
