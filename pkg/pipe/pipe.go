package pipe

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type TryType struct {
	Type string `yaml:"type"`
}

type Pipe interface {
	Init(node *yaml.Node)
	Handle()
}

func createPipe(pipeType string, node *yaml.Node) Pipe {
	var pipe Pipe

	switch pipeType {
	case "http":
		pipe = &Http{
			Host: "0.0.0.0",
		}
	default:
		log.Fatal("wrong pipe type: %s", pipeType)
	}

	pipe.Init(node)
	return pipe
}

func runPipes(node *yaml.Node) []Pipe {
	var tryTypes []TryType
	var pipes []Pipe

	err := node.Decode(&tryTypes)
	if err != nil {
		log.Fatal("error parsing pipes: %s", err)
	}

	for i, tryType := range tryTypes {
		pipes = append(pipes, createPipe(tryType.Type, node.Content[i]))
	}

	return pipes
}

func New(cfg *config.Config) {
	if cfg.Pipes.Kind != yaml.SequenceNode {
		log.Fatal("pipes block must be an array")
	}

	runPipes(&cfg.Pipes)
}
