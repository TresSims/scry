package main

import "fmt"

type NotABundle struct {
	UselessData string
}

var Bundle NotABundle = NotABundle{}

func main() {
	fmt.Print("That's all folks.")
}
