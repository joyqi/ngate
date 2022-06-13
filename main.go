package main

import (
	"fmt"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/pipe"
)

func main() {
	args := NewArgs()
	cfg := config.New(args.configFile)
	pipe.New(cfg)

	fmt.Println(cfg.Pipes[0])
}
