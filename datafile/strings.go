package datafile

import (
	"bufio"
	"os"
)

// getstrings reads a file and returns its contents as a slice of strings
func GetStrings(filename string) ([]string, error) {
	var lines []string
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	if scanner.Err() != nil {
		return nil, err
	}
	return lines, nil
}
