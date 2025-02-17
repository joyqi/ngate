package auth

import (
	"encoding/json"
	"github.com/joyqi/go-lafi/api/authen"
	"github.com/joyqi/go-lafi/api/contact"
	feishu "github.com/joyqi/go-lafi/oauth2"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
)

type Feishu struct {
	BaseAuth
	conf *feishu.Config
}

func NewFeishu(cfg *config.AuthConfig, url *url.URL) *Feishu {
	kind := feishu.TypeFeishu

	if cfg.Kind == "lark" {
		kind = feishu.TypeLark
	}

	conf := &feishu.Config{
		Type:        kind,
		AppID:       cfg.AppId,
		AppSecret:   cfg.AppSecret,
		RedirectURL: cfg.RedirectURL,
	}

	return &Feishu{NewBaseAuth(url), conf}
}

func (f *Feishu) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	if f.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		token, err := f.conf.Exchange(ctx, string(ctx.QueryArgs().Peek("code")))

		if err != nil {
			ctx.Error("Error Access Token", fasthttp.StatusForbidden)
		} else {
			state := ctx.QueryArgs().Peek("state")

			if len(state) > 0 {
				redirect(string(state))
			} else {
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}

			// save token to session
			f.saveTokenSession(ctx, token, session)
		}
	} else {
		redirect(f.conf.AuthCodeURL(f.RequestURL(ctx)))
	}
}

func (f *Feishu) Valid(ctx *fasthttp.RequestCtx, session Session) bool {
	if token := session.Get("token"); token != "" {
		tk := &feishu.Token{}
		err := json.Unmarshal([]byte(token), tk)

		if err != nil {
			log.Error("token: %s", err)
			return false
		}

		if !tk.Valid() {
			ts := f.conf.TokenSource(ctx, tk)
			newTk, tkErr := ts.Token()

			if tkErr != nil {
				log.Error("token: %s", tkErr)
				return false
			}

			f.saveTokenSession(ctx, newTk, session)
			log.Debug("access_token: refresh")
			return true
		} else {
			log.Debug("access_token: valid")
			return true
		}
	}

	log.Debug("access_token: null")
	return false
}

func (f *Feishu) Groups(ctx *fasthttp.RequestCtx, session Session) []string {
	return strings.Split(session.Get("group"), ",")
}

func (f *Feishu) saveTokenSession(ctx *fasthttp.RequestCtx, token *feishu.Token, session Session) {
	v, _ := json.Marshal(token)
	session.Set("token", string(v))

	// get user info
	user := f.retrieveUserInfo(ctx, token)

	// save group to session
	group := f.retrieveUserGroup(ctx, user.OpenId)
	session.Set("group", group)
}

func (f *Feishu) retrieveUserInfo(ctx *fasthttp.RequestCtx, token *feishu.Token) *authen.UserInfoData {
	client := f.conf.TokenSource(ctx, token).Client()
	api := &authen.UserInfo{Client: client}

	user, err := api.Get()
	if err != nil {
		return nil
	}

	return user
}

func (f *Feishu) retrieveUserGroup(ctx *fasthttp.RequestCtx, openId string) string {
	client := f.conf.TenantTokenSource(ctx).Client()
	api := &contact.Group{Client: client}

	group, err := api.MemberBelong(&contact.GroupMemberBelongParams{MemberId: openId})

	if err != nil {
		return ""
	}

	return strings.Join(group.GroupList, ",")
}
