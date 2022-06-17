package main

import (
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/pkg/auth"
	"github.com/joyqi/ngate/pkg/pipe"
)

var cfg *config.Config

func init() {
	args := NewArgs()
	cfg = config.New(args.configFile)
}

func main() {
	pipe.New(cfg, auth.New(&cfg.Auth))
}
