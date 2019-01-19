package raft

import (
	"log"
	"math/rand"
)

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 1 {
		log.Printf(format, a...)
	}
	return
}
func ADPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}
func randInt(low, upper int) int {
	randNum := rand.Intn(upper - low) + low
	return randNum
}

