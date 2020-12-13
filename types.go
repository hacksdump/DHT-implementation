package main

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

type RoutingInfo struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}