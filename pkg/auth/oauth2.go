package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	Host   string
	Path   string
	Config config.AuthConfig
}

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx) bool {
	if string(ctx.Host()) == oauth.Host && string(ctx.Path()) == oauth.Path {
		ctx.SetBody([]byte("hello"))
		return true
	}

	oauth.auth(ctx)
	return false
}

func (oauth *OAuth2) auth(ctx *fasthttp.RequestCtx) {
	conf := &oauth2.Config{
		ClientID:     oauth.Config.AppKey,
		ClientSecret: oauth.Config.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth.Config.AuthorizeUrl,
			TokenURL: oauth.Config.AccessTokenUrl,
		},
		RedirectURL: oauth.Config.RedirectUrl,
		Scopes:      oauth.Config.Scopes,
	}

	if ctx.QueryArgs().Has("code") {
	}

	ctx.Redirect(conf.AuthCodeURL("state", oauth2.SetAuthURLParam("app_id", oauth.Config.AppKey)), fasthttp.StatusFound)
}
