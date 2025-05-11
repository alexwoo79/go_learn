package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func average(numbers ...float64) float64 {
	var sum float64 = 0
	for _, number := range numbers {
		sum += number
	}
	sampleCount := float64(len(numbers))
	return sum / sampleCount
}

// main function

func main() {
	arguments := os.Args[1:]
	var numbers []float64
	for _, arg := range arguments {
		number, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, number)
	}
	fmt.Println("输入的数据为:", os.Args[1:])
	fmt.Println("输入的数据个数:", len(os.Args[1:]))
	fmt.Printf("Average: %0.2f\n", average(numbers...))
}
