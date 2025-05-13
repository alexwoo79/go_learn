package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// ReadFile reads a file and returns its contents as a slice of float64
func main() {
	// Open the file
	file, err := os.Open("datafile/data.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read each line of the file
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	err = file.Close()

	if err != nil {
		log.Fatal(err)
	}
	// Check for errors during scanning
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
