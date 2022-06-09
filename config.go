package main

import (
	"fmt"
	"github.com/joyqi/dahuang/auth"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	log  *log.Logger
	Data ConfigData
	Auth *auth.Auth
}

type TryAuthConfig struct {
	Type string `yaml:"type"`
}

type ConfigData struct {
	Host string    `yaml:"host"`
	Port int       `yaml:"port"`
	Auth yaml.Node `yaml:"auth"`
}

func NewConfig(file string) *Config {
	c := Config{
		log: log.New(os.Stderr, "dahuang: ", log.Ldate|log.Ltime),
		Data: ConfigData{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}

	c.parse(file)
	c.parseAuth()

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
}

func (c *Config) parseAuth() {
	tryAuthConfig := TryAuthConfig{}
	err := c.Data.Auth.Encode(&tryAuthConfig)

	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing auth config: %s", err))
	}

	var a *auth.Auth

	switch tryAuthConfig.Type {
	case "oauth":
		a = &auth.OAuth{}
	default:
		c.log.Fatal("wrong auth type: %s", tryAuthConfig.Type)
	}

	err = c.Data.Auth.Encode(a)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing %s config: %s", tryAuthConfig.Type, err))
	}

	c.Auth = a
}
