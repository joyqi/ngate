package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Auth  AuthConfig   `yaml:"auth"`
	Pipes []PipeConfig `yaml:"pipes,flow"`
}

type AuthConfig struct {
	Kind string `yaml:"kind"`

	// for oauth
	AppKey         string `yaml:"app_key"`
	AppSecret      string `yaml:"app_secret"`
	AccessTokenUrl string `yaml:"access_token_url"`
	AuthorizeUrl   string `yaml:"authorize_url"`
}

type PipeConfig struct {
	Kind     string          `yaml:"kind"`
	Host     string          `yaml:"host"`
	Port     int             `yaml:"port"`
	Backends []BackendConfig `yaml:"backends,flow"`
}

type BackendConfig struct {
	Kind string `yaml:"kind"`

	// for proxy
	Hostname   string `yaml:"hostname"`
	RemoteAddr string `yaml:"remote_addr"`
	RemotePort int    `yaml:"remote_port"`
}

// New read and parse a yaml file
func New(file string) *Config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("error reading %s: %s", file, err)
	}

	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("error parsing: %s", err)
	}

	return &cfg
}
