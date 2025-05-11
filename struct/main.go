// main.go related with magazine.Subscriber struct
// magazine.Subscriber struct is defined in magazine/magazine.go
package main

import (
	"fmt"
	// import magazine.Subscriber struct
	"github.com/alexwoo79/go_learn/magazine"
)

// struct method
func printInfo(s *magazine.Subscriber) {
	fmt.Printf("Name: %s, Rate: %.2f, Active: %t\n", s.Name, s.Rate, s.Active)
}

// struct method
func defaultSubscriber(inputName string) *magazine.Subscriber {
	var s magazine.Subscriber
	s.Name = inputName
	s.Rate = 5.99
	s.Active = true

	return &s
}

func applyDiscount(s *magazine.Subscriber) {
	s.Rate = 4.99
}

// main function
func main() {
	subscriber := magazine.Subscriber{
		Name: "John Doe"}
	subscriber.City = "New York"
	subscriber.State = "NY"
	subscriber.Street = "123 Main St"
	subscriber.ZipCode = "10001"
	fmt.Println("Subscriber:", subscriber.Name)
	fmt.Println("Street:", subscriber.Street)
	fmt.Println("City:", subscriber.City)
	fmt.Println("State:", subscriber.State)
	fmt.Println("ZipCode:", subscriber.ZipCode)
	employee := magazine.Employee{Name: "Alex Woo"}
	employee.City = "中国北京"
	employee.State = "门头沟"
	employee.Street = "金安路"
	employee.ZipCode = "1000111"
	fmt.Println("Subscriber:", employee.Name)
	fmt.Println("Street:", employee.Street)
	fmt.Println("City:", employee.City)
	fmt.Println("State:", employee.State)
	fmt.Println("ZipCode:", employee.ZipCode)
}
