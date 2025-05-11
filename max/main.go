package main

import (
	"fmt"
	"math"
)

func max(numbers ...float64) float64 {
	max_value := math.Inf(-1)
	for _, number := range numbers {
		if number > max_value {
			max_value = number
		}
	}
	return max_value

}

func main() {
	// Create an array of integers
	numbers := []float64{1.0, 7.0, 3.0, 4.0, 5.0}
	max_value := max(numbers...)
	fmt.Println(max_value)
}
