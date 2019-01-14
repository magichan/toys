package raft

import (
	"log"
	"math/rand"
)

// Debugging
const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}
func randInt(low, upper int) int {
	return low + rand.Int()*(upper-low)
}

