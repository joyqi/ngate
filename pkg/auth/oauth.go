package auth

import (
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type OAuth struct {
	AppKey         string `yaml:"app_key"`
	AppSecret      string `yaml:"app_secret"`
	AccessTokenUrl string `yaml:"access_token_url"`
	AuthorizeUrl   string `yaml:"authorize_url"`
}

func (oauth *OAuth) missingField() string {
	switch 0 {
	case len(oauth.AppKey):
		return "app_key"
	case len(oauth.AppSecret):
		return "app_secret"
	case len(oauth.AccessTokenUrl):
		return "access_token_url"
	case len(oauth.AuthorizeUrl):
		return "authorize_url"
	default:
		return ""
	}
}

func (oauth *OAuth) Init(node *yaml.Node) {
	err := node.Decode(oauth)
	if err != nil {
		log.Fatal("error parsing oauth config: %s", err)
	}

	missingFiled := oauth.missingField()

	if len(missingFiled) > 0 {
		log.Fatal("missing auth filed: %s", missingFiled)
	}
}

func (oauth *OAuth) Handler() {

}
