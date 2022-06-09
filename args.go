package main

import "flag"

type Args struct {
	configFile string
}

func NewArgs() *Args {
	args := Args{}

	flag.StringVar(&args.configFile, "config", "config.yaml", "Config file path")
	flag.Parse()

	return &args
}
