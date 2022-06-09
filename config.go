package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	log  *log.Logger
	Data ConfigData
	AuthType
}

type AuthType struct {
	Type string `yaml:"type"`
}

type ConfigData struct {
	Host string    `yaml:"host"`
	Port int       `yaml:"port"`
	Auth yaml.Node `yaml:"auth"`
}

// NewConfig read and parse a yaml file
func NewConfig(file string) *Config {
	c := Config{
		log: log.New(os.Stderr, "dahuang: ", log.Ldate|log.Ltime),
		Data: ConfigData{
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
		c.log.Fatal(fmt.Sprintf("error reading %s: %s", file, err))
	}

	err = yaml.Unmarshal(data, &c.Data)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing config: %s", err))
	}

	err = c.Data.Auth.Decode(&c.AuthType)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing auth config: %s", err))
	}
}
