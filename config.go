package main

import "fmt"

const K uint = 3
const BasePort uint = 8000

func CompleteBaseAddress() string {
	return fmt.Sprintf(":%d", BasePort)
}
