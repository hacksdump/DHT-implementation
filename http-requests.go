package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func readDataFromSpecificNode(node Node, key string) (KeyValuePair, error) {
	url := fmt.Sprintf("http://%s/server-data", node.Address)
	jsonData, _ := json.Marshal(KeyValuePair{key, ""})
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return KeyValuePair{}, errors.New(resp.Status)
	}
	var keyValuePair KeyValuePair
	json.NewDecoder(resp.Body).Decode(&keyValuePair)
	return keyValuePair, nil
}

func writeDataToSpecificNode(node Node, keyValuePair KeyValuePair) (ResponseMessage, error) {
	url := fmt.Sprintf("http://%s/server-data", node.Address)
	jsonData, _ := json.Marshal(keyValuePair)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	json.NewDecoder(resp.Body).Decode(&keyValuePair)
	return ResponseMessage{true, "Written"}, nil
}

func deleteDataFromSpecificNode(node Node, key string) (ResponseMessage, error) {
	url := fmt.Sprintf("http://%s/server-data", node.Address)
	jsonData, _ := json.Marshal(KeyValuePair{key, ""})
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		return ResponseMessage{true, "Successfully deleted"}, nil
	}
	return ResponseMessage{false, "Delete failure"}, errors.New("delete failure")
}

func getAllDataFromNode(node Node) []KeyValuePair {
	resp, err := http.Get(fmt.Sprintf("http://%s/data", node.Address))
	if err == nil {
		var allData []KeyValuePair
		json.NewDecoder(resp.Body).Decode(&allData)
		return allData
	}
	return []KeyValuePair{}
}