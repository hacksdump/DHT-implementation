package main

import (
	"encoding/json"
	"log"
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
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
