package main

import "fmt"

// Greeting returns the classic greeting printed by jcli-hello.
func Greeting() string {
	return "Hello World"
}

func main() {
	fmt.Println(Greeting())
}
