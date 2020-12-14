package main

import (
	"DHT-implementation/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var node Node

func getNextNode(node Node) Node {
	nextNode, err := addrToNode(node.NextAddress)
	if err != nil {
		log.Fatal("Get next node failed")
	}
	return nextNode
}

func findNode(start Node, key string) Node {
	current := start
	keyHash := util.Hash(key, K)
	for util.RingDistance(util.Hash(current.Address, K), keyHash, K) >
		util.RingDistance(util.Hash(current.Next().Address, K), keyHash, K) {
		current = getNextNode(current)
	}
	return current
}

func insertSelfAfter(oldNodeAddr string) {
	prevNode, err := addrToNode(oldNodeAddr)
	if err != nil {
		log.Fatal("Insert self in network failed")
	}
	nextNode := prevNode.Next()
	updateRouting(prevNode.Address, RoutingInfo{Next: node.Address})
	updateRouting(nextNode.Address, RoutingInfo{Prev: node.Address})
	node.PrevAddress = prevNode.Address
	node.NextAddress = nextNode.Address
	log.Printf("Node info: %+v", node)
}

func setRoutes() {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/node", getNodeInfoHandler)
	http.HandleFunc("/data", dataDisplayHandler)
	http.HandleFunc("/update-routing", updateRoutingHandler)
	http.HandleFunc("/server-data", directDataLevelCommunicationHandler)
}

func initNode() {
	if util.CheckPortOpen(BasePort) {
		log.Printf("Starting first node at port %d", BasePort)
		node.ID = util.Hash(CompleteBaseAddress(), K)
		node.Address = CompleteBaseAddress()
		node.NextAddress = CompleteBaseAddress()
		node.PrevAddress = CompleteBaseAddress()
		log.Printf("Base node info: %+v", node)
		err := http.ListenAndServe(CompleteBaseAddress(), nil)
		if err != nil {
			log.Fatal("Could not start root node")
		}
	} else {
		// Finding next available port
		availablePort := BasePort + 1
		const MaxTCPPort = 65535
		for !util.CheckPortOpen(availablePort) && availablePort < MaxTCPPort {
			availablePort += 1
		}
		log.Printf("Starting node at port %d", availablePort)
		completeAddress := fmt.Sprintf(":%d", availablePort)
		node.ID = util.Hash(completeAddress, K)
		node.Address = completeAddress
		rootNode, _ := addrToNode(CompleteBaseAddress())
		nodeToInsertAfter := findNode(rootNode, completeAddress)
		insertSelfAfter(nodeToInsertAfter.Address)
		err := http.ListenAndServe(completeAddress, nil)
		if err != nil {
			log.Fatal("Could not start node")
		}
	}
}

func handleRequests() {
	setRoutes()
	initNode()
}

func closeGracefully() {
	prev := node.Prev()
	next := node.Next()
	updateRouting(prev.Address, RoutingInfo{Next: next.Address})
	updateRouting(next.Address, RoutingInfo{Prev: prev.Address})
}

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		closeGracefully()
		os.Exit(0)
	}()
}

func main() {
	setupCloseHandler()
	handleRequests()
}
