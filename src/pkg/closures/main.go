// https://gobyexample.com/closures
package main

import "fmt"

func GetNameGeneratorWithIncrementalId() func() string {
	var id = 0
	// The returned function closes over the variable id to form a closure.
	return func() string {
		id++
		return fmt.Sprintf("My ID #%d", id)
	}
}

func main() {
	firstGenerateNameFunction := GetNameGeneratorWithIncrementalId()
	secondGenerateNameFunction := GetNameGeneratorWithIncrementalId()

	for i := 0; i < 3; i++ {
		fmt.Printf("firstGenerateNameFunction: %s\n", firstGenerateNameFunction())
	}

	for i := 0; i < 3; i++ {
		fmt.Printf("secondGenerateNameFunction: %s\n", secondGenerateNameFunction())
	}
}
