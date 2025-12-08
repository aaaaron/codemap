// Test Go file for parser testing
package main

import "fmt"

// Greet greets a person by name
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// Calculator represents a simple calculator
type Calculator struct {
	result int
}

// NewCalculator creates a new calculator
func NewCalculator() *Calculator {
	return &Calculator{result: 0}
}

// Add adds two numbers
func (c *Calculator) Add(x, y int) int {
	return x + y
}

// Multiply multiplies two numbers
func (c *Calculator) Multiply(x, y int) int {
	return x * y
}

func main() {
	calc := NewCalculator()
	fmt.Println(Greet("World"))
	fmt.Println(calc.Add(2, 3))
}