package main

var keyValueStore = make(map[string]string)

func getAllData() []KeyValuePair {
	allValues := make([]KeyValuePair, 0, len(keyValueStore))
	for k := range keyValueStore {
		allValues = append(allValues, KeyValuePair{k, keyValueStore[k]})
	}
	return allValues
}
