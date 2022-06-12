package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type Pipe struct {
	Type string `yaml:"type"`
}

type PipeConfig interface {
	Validate()
}

type PipeParser struct {
	Backends yaml.Node `yaml:"backends,flow"`
}

type HttpPipe struct {
	Pipe
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (p *HttpPipe) Validate() {
	if p.Port <= 0 {
		log.Fatal("wrong http port number: %d", p.Port)
	}
}

func NewPipes(node *yaml.Node) *[]PipeConfig {
	var tryPipes []Pipe
	var pipes []PipeConfig

	err := node.Decode(&tryPipes)
	if err != nil {
		log.Fatal("error parsing pipes: %s", err)
	}

	for i, tryPipe := range tryPipes {
		pipes = append(pipes, createPipe(tryPipe.Type, node.Content[i]))
	}

	return &pipes
}

func createPipe(pipeType string, node *yaml.Node) PipeConfig {
	var pipe PipeConfig

	switch pipeType {
	case "http":
		pipe = &HttpPipe{
			Host: "0.0.0.0",
		}
	default:
		log.Fatal("wrong pipe type: %s", pipeType)
	}

	err := node.Decode(pipe)
	if err != nil {
		log.Fatal("error parsing pipe block: %s", err)
	}

	pipe.Validate()
	return pipe
}
