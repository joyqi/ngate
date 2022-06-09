package main

import "fmt"

func main() {
	args := NewArgs()
	config := NewConfig(args.configFile)
	auth := config.Auth()

	fmt.Println(auth)
}
