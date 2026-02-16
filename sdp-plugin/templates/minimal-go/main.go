// Package main demonstrates a minimal SDP-compatible Go project.
package main

import "fmt"

func main() {
	message := greet("World")
	fmt.Println(message)
}

// greet returns a greeting message for the given name.
func greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
