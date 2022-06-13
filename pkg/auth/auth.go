package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
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

	// get auth type
	return tryType.Type
}

// New parse the auth block of the config file
func New(cfg *config.Config) Auth {
	var a Auth
	return a
}
