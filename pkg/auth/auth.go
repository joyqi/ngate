package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/valyala/fasthttp"
)

type Auth interface {
	Handler(ctx *fasthttp.RequestCtx)
}

// New parse the auth block of the config file
func New(cfg *config.Config) Auth {
	var a Auth

	switch cfg.Auth.Kind {
	case "oauth":
	default:
		a = &OAuth{
			AppKey:         cfg.Auth.AppKey,
			AppSecret:      cfg.Auth.AppSecret,
			AuthorizeUrl:   cfg.Auth.AuthorizeUrl,
			AccessTokenUrl: cfg.Auth.AccessTokenUrl,
		}
	}

	return a
}
