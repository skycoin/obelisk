package main

type NodeMessage struct {
	arrivalTimeTick int
	sentTimeTick    int
	from    *Node
	to    *Node
	message 	string
}
