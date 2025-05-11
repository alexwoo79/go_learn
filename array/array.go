package main

import (
	"fmt"
)

// Array is a struct that holds an array of integers
func main() {
	// Create an array of integers
	notes := [7]string{"A", "B", "C", "D", "E", "F", "G"}
	// for i := 0; i < len(notes); i++ {
	// 	fmt.Println("Note", i+1, "is", notes[i])
	// }
	for i, note := range notes {
		fmt.Println("Note", i+1, "is", note)
	}
}
