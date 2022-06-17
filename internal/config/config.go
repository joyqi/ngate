package config

import (
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
	AppId          string   `yaml:"app_id"`
	ClientId       string   `yaml:"client_id"`
	AppSecret      string   `yaml:"app_secret"`
	AccessTokenURL string   `yaml:"access_token_url"`
	AuthorizeURL   string   `yaml:"authorize_url"`
	RedirectURL    string   `yaml:"redirect_url"`
	Scopes         []string `yaml:"scopes,flow"`
}

type PipeConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Session SessionConfig `yaml:"session"`
	Backend BackendConfig `yaml:"backend"`
}

type SessionConfig struct {
	CookieKey    string `yaml:"cookie_key"`
	CookieDomain string `yaml:"cookie_domain"`
	HashKey      string `yaml:"hash_key"`
	BlockKey     string `yaml:"block_key"`
	ExpiresIn    int64  `yaml:"expires_in"`
}

type BackendConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int64  `yaml:"timeout"`
}

// New read and parse a yaml file
func New(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
