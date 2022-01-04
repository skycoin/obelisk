package main

import (
	"sort"
)

type NodeMessageQueue struct {
	messages []*NodeMessage
}

func (nmq *NodeMessageQueue) Push(nodeMessage *NodeMessage) {
	nmq.messages = append(nmq.messages, nodeMessage);
	sort.SliceStable(nmq.messages, func(i, j int) bool {
		return nmq.messages[i].arrivalTimeTick < nmq.messages[j].arrivalTimeTick
	})
}

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
