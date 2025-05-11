package main

// Mustang struct
type car struct {
	name     string
	topSpeed int
}

func nitroBoost(c *car) {
	c.topSpeed += 50
}
func main() {
	// Create a new car
	var mustang car
	mustang.name = "Mustang"
	mustang.topSpeed = 225
	// Print the car's name and top speed
	println("Car name:", mustang.name)
	println("Top speed:", mustang.topSpeed)
	// Apply nitro boost
	nitroBoost(&mustang)
	// Print the new top speed
	println("After nitroBoost New top speed:", mustang.topSpeed)
}
