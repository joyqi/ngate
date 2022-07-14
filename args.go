package main

import "flag"

type Args struct {
	ConfigFile string
	DebugMode  bool
}

func NewArgs() *Args {
	args := Args{}

	flag.StringVar(&args.ConfigFile, "c", "config.yaml", "Config file path")
	flag.BoolVar(&args.DebugMode, "d", false, "Debug mode")
	flag.Parse()

	return &args
}
