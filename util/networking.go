package util

import (
	"fmt"
	"net"
)

func CheckPortOpen(port uint) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d",port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}