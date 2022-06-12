package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Parser struct {
	Pipes yaml.Node `yaml:",flow"`
	Auth  yaml.Node `yaml:"auth"`
}

type Config struct {
	Auth  AuthConfig
	Pipes []PipeConfig
}

func (p *Parser) parse(file string) *Config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("error reading %s: %s", file, err)
	}

	err = yaml.Unmarshal(data, &p)
	if err != nil {
		log.Fatal("error parsing: %s", err)
	}

	return &Config{
		Auth:  NewAuth(&p.Auth),
		Pipes: NewPipes(&p.Pipes),
	}
}

// New read and parse a yaml file
func New(file string) *Config {
	p := Parser{}
	return p.parse(file)
}
