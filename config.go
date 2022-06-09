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
	Data *ConfigData
}

type ConfigData struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig(file string) *Config {
	c := Config{
		log: log.New(os.Stderr, "dahuang: ", log.Ldate|log.Ltime),
		Data: &ConfigData{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}

	c.read(file)

	return &c
}

func (c *Config) read(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error reading %s: %s", file, err))
	}

	err = yaml.Unmarshal(data, c.Data)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing config: %s", err))
	}
}
