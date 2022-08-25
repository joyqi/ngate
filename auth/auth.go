package auth

import (
	"context"
	"fmt"
	"github.com/joyqi/ngate/auth/adapters"
	"github.com/joyqi/ngate/auth/session"
	"github.com/joyqi/ngate/config"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Auth interface {
	Init(config *config.AuthConfig) error
	RetrieveToken(ctx *fasthttp.RequestCtx) (*session.Token, error)
	ValidURL(host string, path string) bool
	ValidToken(ctx context.Context, token *session.Token) (*session.Token, error)
	AuthURL(url string) string
}

// New parse the auth block of the config file
func New(cfg *config.AuthConfig) (Auth, error) {
	var a Auth

	switch cfg.Kind {
	case "fake":
	case "feishu":
		a = &adapters.Feishu{}
	case "oauth2":
		fallthrough
	default:
		return nil, fmt.Errorf("wront auth kind: %s", cfg.Kind)
	}

	err := a.Init(cfg)
	if err != nil {
		return nil, err
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
