package main

import "log"

type Node struct {
	ID          uint
	Address     string
	NextAddress string
	PrevAddress string
}

func (node Node) Next() Node {
	nextNode, err := addrToNode(node.NextAddress)
	if err == nil {
		return nextNode
	}
	log.Fatal("Couldn't find prev node")
	return Node{}
}

func (node Node) Prev() Node {
	prevNode, err := addrToNode(node.PrevAddress)
	if err == nil {
		return prevNode
	}
	log.Fatal("Couldn't find next node")
	return Node{}
}