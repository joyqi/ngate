package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Auth interface {
	Handler(ctx *fasthttp.RequestCtx) bool
}

// New parse the auth block of the config file
func New(cfg *config.Config) Auth {
	var a Auth

	u, err := url.Parse(cfg.Auth.RedirectUrl)
	if err != nil {
		log.Fatal("wrong redirect url: %s", cfg.Auth.RedirectUrl)
	}

	switch cfg.Auth.Kind {
	case "oauth2":
		fallthrough
	default:
		a = &OAuth2{
			u.Host,
			u.Path,
			cfg.Auth,
		}
	}

	return a
}
