package main

import (
	"github.com/joyqi/ngate/auth"
	"github.com/joyqi/ngate/config"
	"github.com/joyqi/ngate/server"
	log "github.com/sirupsen/logrus"
	"os"
)

var cfg *config.Config

func init() {
	var err error

	args := NewArgs()

	log.SetOutput(os.Stdout)

	if args.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	if cfg, err = config.ReadFile(args.ConfigFile); err != nil {
		log.Fatal("error parsing config: %s", err)
	}
}

func main() {
	a, err := auth.New(&cfg.Auth)
	if err != nil {
		log.Fatal("error create auth: %s", err)
	}

	if err = server.New(cfg, a); err != nil {
		log.Fatal("error create server: %s", err)
	}
}
