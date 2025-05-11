package datafile

import (
	"bufio"
	"os"
	"strconv"
)

// ReadFile reads a file and returns its contents as a slice of float64
func GetFloats(filename string) ([]float64, error) {
	var numbers []float64
	var number float64
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number, err = strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	if scanner.Err() != nil {
		return nil, err
	}
	return numbers, nil
}
