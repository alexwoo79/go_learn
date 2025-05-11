package main

import (
	"fmt"
	"sort"
)

func main() {
	grades := map[string]float64{
		"John":  90.5,
		"Jane":  85.0,
		"Bob":   78.0,
		"Alice": 92.0,
	}
	/* 构造names列表 */
	var names []string
	for name := range grades {
		names = append(names, name)
	}
	sort.Strings(names) // Sort the names alphabetically
	// Iterate over the map using a for range loop
	for _, name := range names {
		// Print the name and grade
		fmt.Printf("%s: %.2f\n", name, grades[name])
	}
}
