package main

type Node struct {
	ID          uint
	Address     string
	NextAddress string
	PrevAddress string
}

func (node Node) Next() (Node, error) {
	nextNode, err := addrToNode(node.NextAddress)
	if err == nil {
		return nextNode, nil
	}
	return Node{}, err
}

func (node Node) Prev() (Node, error) {
	prevNode, err := addrToNode(node.PrevAddress)
	if err == nil {
		return prevNode, nil
	}
	return Node{}, err
}