package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/logger"
)

type Auth interface {
	New()
}

// NewAuth parse the auth block of the config file
func NewAuth(c *config.Config) Auth {
	err := c.Data.Auth.Decode(&c.AuthType)

	if err != nil {
		logger.Fatal("error parsing auth config: %s", err)
	}

	var a Auth

	switch c.AuthType.Type {
	case "oauth":
		a = &OAuth{}
	default:
		logger.Fatal("wrong auth type: %s", c.AuthType.Type)
	}

	err = c.Data.Auth.Decode(a)
	if err != nil {
		logger.Fatal("error parsing %s config: %s", c.AuthType.Type, err)
	}

	logger.Success("using auth: %s", c.AuthType.Type)
	return a
}
