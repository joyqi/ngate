package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Auth interface {
	Handler(ctx *fasthttp.RequestCtx) string
}

// New parse the auth block of the config file
func New(cfg *config.AuthConfig) Auth {
	var a Auth

	u, err := url.Parse(cfg.RedirectURL)
	if err != nil {
		log.Fatal("wrong redirect url: %s", cfg.RedirectURL)
	}

	switch cfg.Kind {
	case "feishu":
		a = &Feishu{BaseAuth{u}, cfg, 0, nil}
	case "oauth2":
		fallthrough
	default:
		a = &OAuth2{BaseAuth{u}, cfg}
	}

	return a
}

type BaseAuth struct {
	BaseURL *url.URL
}

func (a *BaseAuth) RequestURL(ctx *fasthttp.RequestCtx) string {
	return a.BaseURL.Scheme + "//" + string(ctx.Host()) + string(ctx.RequestURI())
}

func (a *BaseAuth) IsCallback(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Host()) == a.BaseURL.Host && string(ctx.Path()) == a.BaseURL.Path
}
