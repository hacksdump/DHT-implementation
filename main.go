package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

var keyValueStore = make(map[string]string)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Key struct {
	Key string `json:"key"`
}

type Value struct {
	Value string `json:"value"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

const K = 3

func ringDistance(a int, b int) uint {
	if a == b {
		return 0
	}
	if a < b {
		return uint(b - a)
	}
	return uint(int(math.Pow(2, K)) + b - a)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		var requestBody Key
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if val, ok := keyValueStore[requestBody.Key]; ok {
			w.WriteHeader(http.StatusFound)
			err = json.NewEncoder(w).Encode(Value{val})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(ResponseMessage{"Key "  + requestBody.Key + " was not found :("})
		}
	} else if r.Method == http.MethodPut {
		var body KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Key == "" {
			http.Error(w, "No key in request body", http.StatusBadRequest)
		}
		if body.Value == "" {
			http.Error(w, "No value in request body", http.StatusBadRequest)
		}
		keyValueStore[body.Key] = body.Value
		responseMessage := ResponseMessage{"Successfully Written"}
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

func handleRequests() {
	http.HandleFunc("/", homePage)
	const BasePort = 8000
	completeAddress := fmt.Sprintf(":%d", BasePort)
	log.Printf("Trying to start as first node at port %d", BasePort)
	err := http.ListenAndServe(completeAddress, nil)
	currentPort := BasePort + 1
	completeAddress = fmt.Sprintf(":%d", currentPort)
	for err != nil {
		log.Printf("Trying to start as subsequent node at port %d", currentPort)
		err = http.ListenAndServe(completeAddress, nil)
		currentPort += 1
		completeAddress = fmt.Sprintf(":%d", currentPort)
	}
}

func main() {
	handleRequests()
}
