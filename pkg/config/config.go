package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
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
	ExpireHours  int    `yaml:"expire_hours"`
}

type BackendConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// New read and parse a yaml file
func New(file string) *Config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("error reading %s: %s", file, err)
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("error parsing: %s", err)
	}

	return &cfg
}
