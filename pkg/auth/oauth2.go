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

func (oauth *OAuth2) Handler(ctx *fasthttp.RequestCtx) (string, string) {
	conf := oauth.config()
	requestUrl := oauth.RequestUrl(ctx)
	appId := oauth2.SetAuthURLParam("app_id", oauth.Config.AppId)

	if oauth.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		token, err := conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")), appId)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusForbidden)
			return "", ""
		}

		return token.AccessToken, ""
	}

	ctx.Redirect(conf.AuthCodeURL("state", appId), fasthttp.StatusFound)
	return "", requestUrl
}

func (oauth *OAuth2) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauth.Config.ClientId,
		ClientSecret: oauth.Config.AppSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth.Config.AuthorizeUrl,
			TokenURL: oauth.Config.AccessTokenUrl,
		},
		RedirectURL: oauth.Config.RedirectUrl,
		Scopes:      oauth.Config.Scopes,
	}
}
