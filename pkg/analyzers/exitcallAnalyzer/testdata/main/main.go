package main

import "os"

func main() {
	test()
	os.Exit(0) // want "os.Exit in main package"
}

func test() {
	os.Exit(0) // want "os.Exit in main package"
}
