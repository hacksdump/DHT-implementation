package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client side http requests

func addrToNode(addr string) (Node, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/node", addr))
	if err == nil {
		var node Node
		json.NewDecoder(resp.Body).Decode(&node)
		return node, nil
	}
	return Node{}, err
}

func updateRouting(addr string, routingInfo RoutingInfo) (Node, error){
	url := fmt.Sprintf("http://%s/update-routing", addr)
	jsonData, _ := json.Marshal(routingInfo)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var node Node
	json.NewDecoder(resp.Body).Decode(&node)
	return node, nil

	return Node{}, err
}