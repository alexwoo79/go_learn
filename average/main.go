package main

import (
	"fmt"
	"log"

	"github.com/alexwoo79/go_learn/datafile"
)

// Average is a function that calculates the average of an array of integers
func Average(numbers []float64) float64 {
	// Check if the array is empty
	if len(numbers) == 0 {
		return 0
	}
	// Calculate the sum of the numbers
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	// Calculate the average
	average := float64(sum) / float64(len(numbers))
	return average

}

// Variance is a function that calculates the variance of an array of integers
func main() {
	// Create an array of integers
	numbers, err := datafile.GetFloats("/Users/alex/go/src/github.com/alexwoo79/go_learn/datafile/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Calculate the average
	average := Average(numbers)
	fmt.Printf("%#v's Average: %.2f \n", numbers, average)

}
