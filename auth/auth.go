package auth

import (
	"fmt"
	"github.com/joyqi/ngate/config"
	"github.com/valyala/fasthttp"
	"net/url"
)

type SoftRedirect func(url string)
type PipeGroupValid func(group string, hostName string) bool

type Session interface {
	Set(key string, value string)
	SetInt(key string, value int64)
	Get(key string) string
	GetInt(key string) int64
	Delete(key string)
	Expired(last int64) bool
}

type Auth interface {
	Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect)
	Valid(session Session) bool
	GroupValid(hostName string, session Session, valid PipeGroupValid) bool
}

// New parse the auth block of the config file
func New(cfg *config.AuthConfig) (Auth, error) {
	var a Auth

	u, err := url.Parse(cfg.RedirectURL)
	if err != nil {
		return nil, fmt.Errorf("wrong redirect url: %s", cfg.RedirectURL)
	}

	switch cfg.Kind {
	case "fake":
		a = &Fake{}
	case "feishu":
		a = &Feishu{BaseAuth{u}, cfg}
	case "oauth2":
		fallthrough
	default:
		a = &OAuth2{BaseAuth{u}, cfg}
	}

	return a, nil
}

type BaseAuth struct {
	BaseURL *url.URL
}

func (a *BaseAuth) RequestURL(ctx *fasthttp.RequestCtx) string {
	return a.BaseURL.Scheme + "://" + string(ctx.Host()) + string(ctx.RequestURI())
}

func (a *BaseAuth) IsCallback(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Host()) == a.BaseURL.Host && string(ctx.Path()) == a.BaseURL.Path
}
