package mypkg

import "fmt"

// 定义接口
type MyInterface interface {
	MethodWithoutArguments()       //定义接口的方法 1
	MethodWithArguments(float64)   //定义接口的方法 2
	MethodWithReturnValue() string //定义接口的方法 3
}
type MyType int //定义一个类型

// 定义类型实现接口的方法
func (s MyType) MethodWithoutArguments() {
	fmt.Println("MethodWithoutArguments called")
}
func (s MyType) MethodWithArguments(f float64) {
	fmt.Println("MethodWithArguments called with f =", f)
}
func (s MyType) MethodWithReturnValue() string {
	fmt.Println("MethodWithReturnValue")
	return "Hi from MethodWithReturnValue called"
}
func (myType MyType) MethodNotInInterface() {
	fmt.Println("ethodNotInInterface called")
}
