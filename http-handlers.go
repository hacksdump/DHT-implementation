package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		var requestBody KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requestKey := requestBody.Key
		log.Println("Received request to find key", requestKey)
		baseNode, err := addrToNode(CompleteBaseAddress())
		if err != nil {
			log.Fatal("Base node not found while searching for a key")
		}
		readFromNode := findNode(baseNode, requestKey)
		log.Printf("Trying to find the key %s in node %+v", requestKey, readFromNode)
		keyValuePair, err := readDataFromSpecificNode(readFromNode, requestKey)

		w.Header().Set("Content-Type", "application/json")
		if err == nil {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(keyValuePair)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(ResponseMessage{false,"Key " + requestKey + " was not found :("})
		}
	} else if r.Method == http.MethodPut {
		var body KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if body.Key == "" || body.Value == "" {
			responseMessage := ResponseMessage{false,"Key or value missing. No write or update performed"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(responseMessage)
			return
		}

		log.Printf("Received request to save data %+v", body)

		baseNode, err := addrToNode(CompleteBaseAddress())
		if err != nil {
			log.Fatal("Base node not found while searching for a key")
		}
		writeToNode := findNode(baseNode, body.Key)
		_, err = writeDataToSpecificNode(writeToNode, body)
		if err == nil {
			log.Printf("Writing to node %+v", writeToNode)
			responseMessage := ResponseMessage{true,
				fmt.Sprintf("Successfully Written to node with ID: %d and address: %s",
					writeToNode.ID,
					writeToNode.Address)}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(responseMessage)
		} else {
			log.Fatal(err)
		}
	} else {
		http.Error(w, "Only use GET or PUT", http.StatusMethodNotAllowed)
	}
}

func directDataLevelCommunicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestBody KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		requestKey := requestBody.Key
		log.Printf("Checking if key %s is present in local store", requestKey)
		if val, ok := keyValueStore[requestKey]; ok {
			keyValuePair := KeyValuePair{requestKey, val}
			log.Println("Found", keyValuePair)
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(keyValuePair)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("Not found")
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(ResponseMessage{false,"Key " + requestKey + " was not found :("})
		}
	} else if r.Method == http.MethodPut {
		var body KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if body.Key == "" || body.Value == "" {
			responseMessage := ResponseMessage{false,"Key or value missing. No write or update performed"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(responseMessage)
			return
		}
		keyValueStore[body.Key] = body.Value
		log.Printf("Written the data %+v", KeyValuePair{body.Key, body.Value})
		responseMessage := ResponseMessage{true,"Successfully Written"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(responseMessage)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Error(w, "Only use GET or PUT", http.StatusMethodNotAllowed)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Node address: %s", node.Address)
}

func dataDisplayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getAllData())
}

func getNodeInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(node)
	} else {
		http.Error(w, "Only use GET", http.StatusMethodNotAllowed)
	}
}

func updateRoutingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var newRoutingInfo RoutingInfo
		err := json.NewDecoder(r.Body).Decode(&newRoutingInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(newRoutingInfo.Prev) > 0 {
			node.PrevAddress = newRoutingInfo.Prev
		}
		if len(newRoutingInfo.Next) > 0 {
			node.NextAddress = newRoutingInfo.Next
		}
		log.Printf("Updated node info: %+v", node)
	} else {
		http.Error(w, "Only use GET or PUT", http.StatusMethodNotAllowed)
	}
}
