package auth

import "fmt"

type Auth interface {
	New()
}

// NewAuth parse auth block of the config file
func NewAuth() Auth {
	err := c.Data.Auth.Decode(&c.AuthType)

	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing auth config: %s", err))
	}

	var a auth.Auth

	switch tryAuthConfig.Type {
	case "oauth":
		a = &auth.OAuth{}
	default:
		c.log.Fatal("wrong auth type: %s", tryAuthConfig.Type)
	}

	err = c.Data.Auth.Decode(a)
	if err != nil {
		c.log.Fatal(fmt.Sprintf("error parsing %s config: %s", tryAuthConfig.Type, err))
	}

	return a
}
