package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"gopkg.in/yaml.v3"
)

type TryType struct {
	Type string `yaml:"type"`
}

type Auth interface {
	Init(node *yaml.Node)
	Handler()
}

func authType(cfg *config.Config) string {
	tryType := TryType{
		Type: "oauth",
	}

	err := cfg.Auth.Decode(&tryType)
	if err != nil {
		log.Fatal("error parsing auth block: %s", err)
	}

	// get auth type
	return tryType.Type
}

// New parse the auth block of the config file
func New(cfg *config.Config) Auth {
	authType := authType(cfg)
	var a Auth

	switch authType {
	case "oauth":
		a = &OAuth{}
	default:
		log.Fatal("wrong auth type: %s", authType)
	}

	a.Init(&cfg.Auth)
	log.Success("using auth: %s", authType)
	return a
}
