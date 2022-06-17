package main

import (
	"github.com/joyqi/dahuang/internal/config"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/pipe"
)

var cfg *config.Config

func init() {
	args := NewArgs()
	cfg = config.New(args.configFile)
}

func main() {
	pipe.New(cfg, auth.New(&cfg.Auth))
}
