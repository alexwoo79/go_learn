package magazine

// struct
type Subscriber struct {
	Name    string // first character is capitalized, so it is exported
	Rate    float64
	Active  bool
	Address // struct field
}

type Employee struct {
	Name    string
	Salary  float64
	Address // struct field
}
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}
