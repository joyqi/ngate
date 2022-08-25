package adapters

import (
	"context"
	"fmt"
	"github.com/joyqi/ngate/auth/session"
	"github.com/joyqi/ngate/config"
	"github.com/joyqi/ngate/internal/go-oauth2-cn/feishu"
	"github.com/valyala/fasthttp"
	"net/url"
	"time"
)

type Feishu struct {
	conf        *feishu.Config
	redirectURL *url.URL
}

func (f *Feishu) Init(config *config.AuthConfig) error {
	f.conf = &feishu.Config{
		AppID:       config.AppId,
		AppSecret:   config.AppSecret,
		RedirectURL: config.RedirectURL,
	}

	u, err := url.Parse(config.RedirectURL)
	if err != nil {
		return err
	}

	f.redirectURL = u
	return nil
}

func (f *Feishu) ValidURL(host string, path string) bool {
	return host == f.redirectURL.Host && path == f.redirectURL.Path
}

func (f *Feishu) ValidToken(ctx context.Context, token *session.Token) (*session.Token, error) {
	tk := &feishu.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       time.Unix(token.ExpiresAt, 0),
		Groups:       token.Groups,
	}

	ts := f.conf.TokenSource(ctx, tk)

	if tk, err := ts.Token(); err != nil {
		return nil, err
	} else {
		return &session.Token{
			AccessToken:  tk.AccessToken,
			RefreshToken: tk.RefreshToken,
			ExpiresAt:    tk.Expiry.Unix(),
			Groups:       tk.Groups,
		}, nil
	}
}

func (f *Feishu) RetrieveToken(ctx *fasthttp.RequestCtx) (*session.Token, error) {
	if !ctx.QueryArgs().Has("code") {
		return nil, fmt.Errorf("no code found")
	}

	code := string(ctx.QueryArgs().Peek("code"))
	if tk, err := f.conf.Exchange(ctx, code); err != nil {
		return nil, err
	} else {
		return &session.Token{
			AccessToken:  tk.AccessToken,
			RefreshToken: tk.RefreshToken,
			ExpiresAt:    tk.Expiry.Unix(),
			Groups:       tk.Groups,
		}, nil
	}
}

func (f *Feishu) AuthURL(url string) string {
	return f.AuthURL(url)
}
