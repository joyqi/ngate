package auth

import (
	"github.com/joyqi/ngate/internal/config"
	"github.com/valyala/fasthttp"
	"github.com/wenerme/go-wecom/wecom"
	"github.com/wenerme/go-wecom/wecom/auth"
	"net/url"
	"strconv"
	"strings"
)

type WecomAuthCodeURL func(state string) string

type Wecom struct {
	BaseAuth
	client      *wecom.Client
	AuthCodeURL WecomAuthCodeURL
}

type WecomUserInfo struct {
	UserID string
	Groups []string
}

func NewWecom(cfg *config.AuthConfig, url *url.URL) *Wecom {
	agentId, _ := strconv.Atoi(cfg.ClientId)
	store := &wecom.SyncMapStore{}

	client := wecom.NewClient(wecom.Conf{
		CorpID:        cfg.AppId,
		AgentID:       agentId,
		CorpSecret:    cfg.AppSecret,
		TokenProvider: &wecom.TokenCache{Store: store},
	})

	authCodeURL := func(state string) string {
		return auth.BuildOAuth2URL(auth.BuildOAuth2URLOptions{
			AppID:        cfg.AppId,
			RedirectURI:  cfg.RedirectURL,
			ResponseType: "code",
			Scope:        "snsapi_base",
			State:        state,
			AgentID:      cfg.ClientId,
		})
	}

	return &Wecom{NewBaseAuth(url), client, authCodeURL}
}

func (w *Wecom) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	if w.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		userInfo, err := w.retrieveUserInfo(string(ctx.QueryArgs().Peek("code")))

		if err != nil {
			ctx.Error("Error Access Token", fasthttp.StatusForbidden)
		} else {
			state := ctx.QueryArgs().Peek("state")

			if len(state) > 0 {
				redirect(string(state))
			} else {
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}

			session.Set("user_id", userInfo.UserID)
			session.Set("groups", strings.Join(userInfo.Groups, ","))
		}
	} else {
		redirect(w.AuthCodeURL(w.RequestURL(ctx)))
	}
}

func (w *Wecom) Valid(ctx *fasthttp.RequestCtx, session Session) bool {
	return session.Get("user_id") != ""
}

func (w *Wecom) Groups(ctx *fasthttp.RequestCtx, session Session) []string {
	return strings.Split(session.Get("groups"), ",")
}

func (w *Wecom) retrieveUserInfo(code string) (*WecomUserInfo, error) {
	accessToken, err := w.client.AccessToken()

	if err != nil {
		return nil, err
	}

	userInfo, userInfoErr := w.client.GetUserInfo(&wecom.GetUserInfoRequest{
		AccessToken: accessToken,
		Code:        code,
	})

	if userInfoErr != nil {
		return nil, userInfoErr
	}

	user, userErr := w.client.GetUser(&wecom.GetUserRequest{
		UserID: userInfo.UserID,
	})

	if userErr != nil {
		return nil, userErr
	}

	groups := make([]string, 0)

	for _, deptId := range user.Department {
		groups = append(groups, strconv.Itoa(deptId))
	}

	return &WecomUserInfo{userInfo.UserID, groups}, nil
}
