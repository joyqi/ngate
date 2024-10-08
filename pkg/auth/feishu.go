package auth

import (
	"bytes"
	"encoding/json"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Endpoint is feishu's endpoint url
const (
	FeishuAuthURL         = "/authen/v1/index"
	FeishuTokenURL        = "/authen/v1/access_token"
	FeishuRefreshTokenURL = "/authen/v1/refresh_access_token"
	FeishuTenantTokenURL  = "/auth/v3/tenant_access_token/internal"
	FeishuUserGroupURL    = "/contact/v3/group/member_belong"
)

// feishu tenant token
var (
	tenantTokenExpireAt int64
	tenantToken         *FeishuTenantToken
)
var refreshLock sync.Map

type Feishu struct {
	BaseAuth
	Config *config.AuthConfig
}

type FeishuTenantToken struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type FeishuAccessToken struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Name         string `json:"name"`
		EnName       string `json:"en_name"`
		OpenId       string `json:"open_id"`
		ExpiresIn    int64  `json:"expires_in"`
	} `json:"data"`
}

type FeishuUserGroup struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		GroupList []string `json:"group_list,flow"`
		HasMore   bool     `json:"has_more"`
	} `json:"data"`
}

func (f *Feishu) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	if f.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		accessToken := f.retrieveToken(string(ctx.QueryArgs().Peek("code")))

		if accessToken != nil && accessToken.Code == 0 {
			state := ctx.QueryArgs().Peek("state")
			group := f.retrieveUserGroup(accessToken.Data.OpenId)

			if len(state) > 0 {
				redirect(string(state))
			} else {
				ctx.Error("Not Found", fasthttp.StatusNotFound)
			}

			session.Set("group", group)
			session.Set("access_token", accessToken.Data.AccessToken)
			session.Set("refresh_token", accessToken.Data.RefreshToken)
			session.SetInt("valid_at", time.Now().Unix())
		} else {
			ctx.Error("Error Access Token", fasthttp.StatusForbidden)
		}
	} else {
		redirect(f.authCodeURL(f.RequestURL(ctx)))
	}
}

func (f *Feishu) Valid(session Session) bool {
	if accessToken := session.Get("access_token"); accessToken != "" {
		if session.Expired(session.GetInt("valid_at")) {
			if _, exists := refreshLock.LoadOrStore(accessToken, 1); !exists {
				refreshToken := f.retrieveRefreshToken(session.Get("refresh_token"))
				valid := false

				if refreshToken != nil && refreshToken.Code == 0 {
					group := f.retrieveUserGroup(refreshToken.Data.OpenId)

					session.Set("group", group)
					session.Set("access_token", refreshToken.Data.AccessToken)
					session.Set("refresh_token", refreshToken.Data.RefreshToken)
					session.SetInt("valid_at", time.Now().Unix())
					valid = true
				} else {
					session.Delete("access_token")
					session.Delete("refresh_token")
					session.Delete("valid_at")
				}

				log.Debug("access_token: refresh: %s", refreshToken)
				refreshLock.Delete(accessToken)
				return valid
			}

			log.Debug("access_token: locked")
		} else {
			log.Debug("access_token: valid")
			return true
		}
	}

	log.Debug("access_token: null")
	return false
}

func (f *Feishu) GroupValid(hostName string, session Session, valid PipeGroupValid) bool {
	return valid(session.Get("group"), hostName)
}

func (f *Feishu) requestTenantToken() {
	if now := time.Now().Unix(); now > tenantTokenExpireAt {
		tenantToken = f.retrieveTenantToken()

		if tenantToken != nil && tenantToken.Code == 0 {
			tenantTokenExpireAt = now + tenantToken.Expire
			log.Success("%s tenant token: %s", f.Config.Kind, tenantToken.TenantAccessToken)
		}
	}
}

func (f *Feishu) retrieveTenantToken() *FeishuTenantToken {
	data := map[string]string{
		"app_id":     f.Config.AppId,
		"app_secret": f.Config.AppSecret,
	}

	body, err := f.postURL(FeishuTenantTokenURL, data, "")
	if err != nil {
		return nil
	}

	token := FeishuTenantToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil
	}

	return &token
}

func (f *Feishu) retrieveToken(code string) *FeishuAccessToken {
	f.requestTenantToken()
	if tenantToken == nil {
		return nil
	}

	data := map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}

	body, err := f.postURL(FeishuTokenURL, data, tenantToken.TenantAccessToken)
	if err != nil {
		return nil
	}

	token := FeishuAccessToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil
	}

	return &token
}

func (f *Feishu) retrieveRefreshToken(refreshToken string) *FeishuAccessToken {
	f.requestTenantToken()
	if tenantToken == nil {
		return nil
	}

	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	body, err := f.postURL(FeishuRefreshTokenURL, data, tenantToken.TenantAccessToken)
	if err != nil {
		return nil
	}

	token := FeishuAccessToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil
	}

	return &token
}

func (f *Feishu) retrieveUserGroup(openId string) string {
	f.requestTenantToken()
	if tenantToken == nil {
		return ""
	}

	body, err := f.getURL(FeishuUserGroupURL+"?member_id="+openId, tenantToken.TenantAccessToken)
	if err != nil {
		return ""
	}

	group := FeishuUserGroup{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		return ""
	}

	return strings.Join(group.Data.GroupList, ",")
}

func (f *Feishu) postURL(url string, data interface{}, token string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	c := &fasthttp.Client{}

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req.SetRequestURI(f.makeURL(url))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json; charset=utf-8")
	req.SetBody(body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	err = c.Do(req, resp)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (f *Feishu) getURL(url string, token string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	c := &fasthttp.Client{}

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(f.makeURL(url))
	req.Header.SetMethod(fasthttp.MethodGet)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	err := c.Do(req, resp)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (f *Feishu) authCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString(f.makeURL(FeishuAuthURL))
	v := url.Values{
		"response_type": {"code"},
		"app_id":        {f.Config.AppId},
	}
	if f.Config.RedirectURL != "" {
		v.Set("redirect_uri", f.Config.RedirectURL)
	}
	if len(f.Config.Scopes) > 0 {
		v.Set("scope", strings.Join(f.Config.Scopes, " "))
	}
	if state != "" {
		v.Set("state", state)
	}

	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

// makeURL adjust the url for feishu or lark
func (f *Feishu) makeURL(url string) string {
	if f.Config.Kind == "feishu" {
		return "https://open.feishu.cn/open-apis" + url
	} else {
		return "https://open.larksuite.com/open-apis" + url
	}
}
