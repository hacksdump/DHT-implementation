package main

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResponseMessage struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type RoutingInfo struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}