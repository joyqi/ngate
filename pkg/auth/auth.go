package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Auth interface {
	Handler(ctx *fasthttp.RequestCtx) (string, string)
}

// New parse the auth block of the config file
func New(cfg *config.AuthConfig) Auth {
	var a Auth

	u, err := url.Parse(cfg.RedirectUrl)
	if err != nil {
		log.Fatal("wrong redirect url: %s", cfg.RedirectUrl)
	}

	switch cfg.Kind {
	case "oauth2":
		fallthrough
	default:
		a = &OAuth2{BaseAuth{u}, cfg}
	}

	return a
}

type BaseAuth struct {
	BaseUrl *url.URL
}

func (a *BaseAuth) RequestUrl(ctx *fasthttp.RequestCtx) string {
	return a.BaseUrl.Scheme + "//" + string(ctx.Host()) + string(ctx.RequestURI())
}

func (a *BaseAuth) IsCallback(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Host()) == a.BaseUrl.Host && string(ctx.Path()) == a.BaseUrl.Path
}
