package main

import (
	"fmt"
)

func main() {
	var ok bool
	ranks := make(map[string]int)
	var rank int
	ranks["A"] = 1
	rank, ok = ranks["A"]
	fmt.Printf("rank: %d ,ok: %v\n", rank, ok)
	// Deleting a key-value pair
	delete(ranks, "A")
	rank, ok = ranks["A"]
	fmt.Printf("rank: %d ,ok: %v\n", rank, ok)

	isPrime := make(map[int]bool)
	var prime bool
	isPrime[2] = true
	prime, ok = isPrime[2]
	fmt.Printf("isPrime: %v ,ok: %v\n", prime, ok)
	// Deleting a key-value pair
	delete(isPrime, 2)
	prime, ok = isPrime[2]
	fmt.Printf("isPrime: %v ,ok: %v\n", prime, ok)
	// Deleting a key-value pair
	delete(isPrime, 3)
	prime, ok = isPrime[3]
	fmt.Printf("isPrime: %v ,ok: %v\n", prime, ok)
}
