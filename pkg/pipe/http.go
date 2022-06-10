package pipe

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type Http struct {
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	Backends yaml.Node `yaml:"backends,flow"`
}

func (http *Http) Init(node *yaml.Node) {
	err := node.Decode(http)
	if err != nil {
		log.Fatal("error parsing pipe config: %s", err)
	}

	if http.Port <= 0 {
		log.Fatal("wrong http port number: %d", http.Port)
	}

	if http.Backends.Kind != yaml.SequenceNode {
		log.Fatal("missing backends")
	}
}

func (http *Http) Handle() {

}
