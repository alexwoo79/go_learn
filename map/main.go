package main

import (
	"fmt"
)

func main() {
	status("John")
	status("Jane")
	status("Bob")
	status("Alice")
	status("Charlie")
}
func status(name string) {
	grades := map[string]float64{
		"John": 85.5,
		"Jane": 92.0,
		"Bob":  78.0}
	grade, ok := grades[name]

	if !ok {
		fmt.Printf("Student %s not found\n", name)
		return
	} else if grade >= 90 {
		fmt.Printf("Student %s has an A\n", name)
	} else if grade >= 80 {
		fmt.Printf("Student %s has a B\n", name)
	} else if grade >= 70 {
		fmt.Printf("Student %s has a C\n", name)
	} else if grade >= 60 {
		fmt.Printf("Student %s has a D\n", name)
	} else {
		fmt.Printf("Student %s has an F\n", name)
	}
}
