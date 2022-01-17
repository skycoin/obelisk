package main

// NodeMessage: Messages exchanged between nodes
type NodeMessage struct {
	arrivalTimeTick int
	sentTimeTick    int
	from    *Node
	to    *Node
	message 	string
}
