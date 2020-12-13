package main

import (
	"DHT-implementation/util"
	"fmt"
	"log"
	"net/http"
)

var keyValueStore = make(map[string]string)

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
	nextNode := getNextNode(current)
	for util.RingDistance(util.Hash(current.Address, K), util.Hash(key, K), K) >
		util.RingDistance(util.Hash(nextNode.Address, K), util.Hash(key, K), K) {
		current = getNextNode(current)
	}
	return current
}

func insertSelfAfter(oldNodeAddr string) {
	prevNode, err := addrToNode(oldNodeAddr)
	if err != nil {
		log.Fatal("Insert self in network failed")
	}
	nextNode, err := prevNode.Next()
	updateRouting(prevNode.Address, RoutingInfo{Next: node.Address})
	updateRouting(nextNode.Address, RoutingInfo{Prev: node.Address})
	node.PrevAddress = prevNode.Address
	node.NextAddress = nextNode.Address
	log.Println("Node info: ", node)
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
		log.Println("Base node info: ",node)
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

func main() {
	handleRequests()
}
