package main

import (
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/joyqi/ngate/pkg/auth"
	"github.com/joyqi/ngate/pkg/pipe"
)

var cfg *config.Config

func init() {
	var err error

	args := NewArgs()

	if args.DebugMode {
		log.DebugMode = true
	}

	if cfg, err = config.New(args.ConfigFile); err != nil {
		log.Fatal("error parsing config: %s", err)
	}
}

func main() {
	a, err := auth.New(&cfg.Auth)
	if err != nil {
		log.Fatal("error create auth: %s", err)
	}

	if err = pipe.New(cfg, a); err != nil {
		log.Fatal("error create pipe: %s", err)
	}
}
