package auth

import (
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	BaseAuth
	Config *config.AuthConfig
}

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx, redirect SoftRedirect) string {
	conf := oauth.config()

	if oauth.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		token, err := conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")))
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusForbidden)
			return ""
		}

		state := ctx.QueryArgs().Peek("state")

		if len(state) > 0 {
			redirect(string(state))
		} else {
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}

		return token.AccessToken
	}

	redirect(conf.AuthCodeURL(oauth.RequestURL(ctx)))
	return ""
}

func (oauth *OAuth2) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauth.Config.ClientId,
		ClientSecret: oauth.Config.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth.Config.AuthorizeURL,
			TokenURL: oauth.Config.AccessTokenURL,
		},
		RedirectURL: oauth.Config.RedirectURL,
		Scopes:      oauth.Config.Scopes,
	}
}
