package auth

import (
	"github.com/joyqi/ngate/internal/config"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
	"net/url"
)

type OAuth2 struct {
	BaseAuth
	conf *oauth2.Config
}

func NewOauth2(cfg *config.AuthConfig, url *url.URL) *OAuth2 {
	conf := &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthorizeURL,
			TokenURL: cfg.AccessTokenURL,
		},
		RedirectURL: cfg.RedirectURL,
		Scopes:      cfg.Scopes,
	}

	return &OAuth2{NewBaseAuth(url), conf}
}

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	if oauth.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		token, err := oauth.conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")))
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
		redirect(oauth.conf.AuthCodeURL(oauth.RequestURL(ctx)))
	}
}

func (oauth *OAuth2) Valid(ctx *fasthttp.RequestCtx, session Session) bool {
	return session.Get("access_token") != ""
}

func (oauth *OAuth2) Groups(ctx *fasthttp.RequestCtx, session Session) []string {
	return nil
}
