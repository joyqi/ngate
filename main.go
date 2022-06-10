package main

import (
	"fmt"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/config"
)

func main() {
	args := NewArgs()
	cfg := config.NewConfig(args.configFile)
	a := auth.NewAuth(cfg)

	fmt.Println(a)
}
