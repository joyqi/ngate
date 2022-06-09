package main

import "fmt"

func main() {
	args := NewArgs()
	config := NewConfig(args.configFile)

	fmt.Println(config.Data)
}
