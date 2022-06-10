package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Pipes yaml.Node `yaml:",flow"`
	Auth  yaml.Node `yaml:"auth"`
}

// New read and parse a yaml file
func New(file string) *Config {
	c := Config{}

	c.parse(file)
	return &c
}

func (c *Config) parse(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("error reading %s: %s", file, err)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("error parsing: %s", err)
	}
}
