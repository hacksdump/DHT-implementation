package util

import (
	"hash/fnv"
	"math"
)

func Hash(text string, bitSpace uint) uint {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return uint(algorithm.Sum32() % uint32(math.Pow(2, float64(bitSpace))))
}

func RingDistance(a uint, b uint, bitSpace uint) uint {
	if a == b {
		return 0
	}
	if a < b {
		return uint(b - a)
	}
	return uint(math.Pow(2, float64(bitSpace))) + b - a
}