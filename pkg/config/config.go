package config

import (
	"github.com/joyqi/dahuang/pkg/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Data
	AuthType
}

type AuthType struct {
	Type string `yaml:"type"`
}

type Data struct {
	Host string    `yaml:"host"`
	Port int       `yaml:"port"`
	Auth yaml.Node `yaml:"auth"`
}

// NewConfig read and parse a yaml file
func NewConfig(file string) *Config {
	c := Config{
		Data: Data{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}

	c.parse(file)
	return &c
}

func (c *Config) parse(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Fatal("error reading %s: %s", file, err)
	}

	err = yaml.Unmarshal(data, &c.Data)
	if err != nil {
		logger.Fatal("error parsing: %s", err)
	}

	err = c.Data.Auth.Decode(&c.AuthType)
	if err != nil {
		logger.Fatal("error parsing auth block: %s", err)
	}
}
