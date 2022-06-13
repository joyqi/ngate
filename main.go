package main

import (
	"fmt"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/pipe"
)

func main() {
	args := NewArgs()
	cfg := config.New(args.configFile)
	a := auth.New(cfg)
	pipe.New(cfg, a)

	fmt.Println(cfg.Pipes[0])
}
