package auth

import (
	"fmt"
	"github.com/joyqi/ngate/internal/config"
	"github.com/valyala/fasthttp"
	"net/url"
)

type SoftRedirect func(url string)
type PipeGroupValid func(groups []string, hostName string) bool

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
	Valid(ctx *fasthttp.RequestCtx, session Session) bool
	Groups(ctx *fasthttp.RequestCtx, session Session) []string
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
		a = NewFake()
	case "feishu", "lark":
		a = NewFeishu(cfg, u)
	case "wecom":
		a = NewWecom(cfg, u)
	case "oauth2":
		fallthrough
	default:
		a = NewOauth2(cfg, u)
	}

	return a, nil
}

type BaseAuth struct {
	BaseURL *url.URL
}

func NewBaseAuth(u *url.URL) BaseAuth {
	return BaseAuth{u}
}

func (a *BaseAuth) RequestURL(ctx *fasthttp.RequestCtx) string {
	return a.BaseURL.Scheme + "://" + string(ctx.Host()) + string(ctx.RequestURI())
}

func (a *BaseAuth) IsCallback(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Host()) == a.BaseURL.Host && string(ctx.Path()) == a.BaseURL.Path
}
