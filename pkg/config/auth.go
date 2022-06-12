package config

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type Auth struct {
	Type string `yaml:"type"`
}

type AuthConfig interface {
	Validate()
}

type OAuthAuth struct {
	Auth
	AppKey         string `yaml:"app_key"`
	AppSecret      string `yaml:"app_secret"`
	AccessTokenUrl string `yaml:"access_token_url"`
	AuthorizeUrl   string `yaml:"authorize_url"`
}

func (oauth *OAuthAuth) Validate() {
	missing := ""

	switch 0 {
	case len(oauth.AppKey):
		missing = "app_key"
	case len(oauth.AppSecret):
		missing = "app_secret"
	case len(oauth.AccessTokenUrl):
		missing = "access_token_url"
	case len(oauth.AuthorizeUrl):
		missing = "authorize_url"
	}

	if len(missing) > 0 {
		log.Fatal("missing auth filed: %s", missing)
	}
}

func NewAuth(node *yaml.Node) AuthConfig {
	tryAuth := Auth{
		Type: "oauth",
	}

	err := node.Decode(&tryAuth)
	if err != nil {
		log.Fatal("error parsing auth block: %s", err)
	}

	var auth AuthConfig

	switch tryAuth.Type {
	case "oauth":
		auth = &OAuthAuth{}
	default:
		log.Fatal("wrong auth type: %s", tryAuth.Type)
	}

	err = node.Decode(auth)
	if err != nil {
		log.Fatal("error parsing auth block: %s", err)
	}

	auth.Validate()
	return auth
}
