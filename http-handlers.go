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
			err = json.NewEncoder(w).Encode(ResponseMessage{"Key " + requestBody.Key + " was not found :("})
		}
	} else if r.Method == http.MethodPut {
		var body KeyValuePair
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if body.Key == "" || body.Value == "" {
			responseMessage := ResponseMessage{"Key or value missing. No write or update performed"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(responseMessage)
			return
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

func statsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Node address: %s", node.Address)
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
	} else {
		http.Error(w, "Only use GET or PUT", http.StatusMethodNotAllowed)
	}
}
