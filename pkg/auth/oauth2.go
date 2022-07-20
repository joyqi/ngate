package auth

import (
	"github.com/joyqi/ngate/internal/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	BaseAuth
	Config *config.AuthConfig
}

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	conf := oauth.config()

	if oauth.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		token, err := conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")))
		if err != nil {
			ctx.Error("Error Access Token", fasthttp.StatusForbidden)
		} else {
			state := ctx.QueryArgs().Peek("state")

			if len(state) > 0 {
				redirect(string(state))
			} else {
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}

			session.Set("access_token", token.AccessToken)
		}
	} else {
		redirect(conf.AuthCodeURL(oauth.RequestURL(ctx)))
	}
}

func (oauth *OAuth2) Valid(session Session) bool {
	return session.Get("access_token") != ""
}

func (oauth *OAuth2) GroupValid(hostName string, session Session, valid PipeGroupValid) bool {
	return true
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
