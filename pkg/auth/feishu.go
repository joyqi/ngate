package auth

import (
	"bytes"
	"encoding/json"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"time"
)

// Endpoint is feishu's endpoint url
const (
	FeishuAuthURL   = "https://open.feishu.cn/open-apis/authen/v1/index"
	FeishuTokenURL  = "https://open.feishu.cn/open-apis/authen/v1/access_token"
	FeishuTenantURL = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
)

type Feishu struct {
	BaseAuth
	Config              *config.AuthConfig
	TenantTokenExpireAt int64
	TenantToken         *FeishuTenantToken
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
		AccessToken string `json:"access_token"`
		Name        string `json:"name"`
		EnName      string `json:"en_name"`
		OpenId      string `json:"open_id"`
		ExpiresIn   int64  `json:"expires_in"`
	} `json:"data"`
}

func (f *Feishu) Handler(ctx *fasthttp.RequestCtx) string {
	if f.IsCallback(ctx) && ctx.QueryArgs().Has("code") {
		if now := time.Now().Unix(); now > f.TenantTokenExpireAt {
			tenantToken := f.retrieveTenantToken()

			if tenantToken != nil {
				f.TenantToken = tenantToken
				f.TenantTokenExpireAt = now + tenantToken.Expire

				log.Success("feishu tenant token: %s", tenantToken.TenantAccessToken)
			}
		}

		var accessToken *FeishuAccessToken

		if f.TenantToken != nil {
			accessToken = f.retrieveToken(string(ctx.QueryArgs().Peek("code")))

			if accessToken != nil {
				state := ctx.QueryArgs().Peek("state")

				if len(state) > 0 {
					ctx.Redirect(string(state), fasthttp.StatusFound)
				} else {
					ctx.Error("Not Found", fasthttp.StatusNotFound)
				}

				return accessToken.Data.AccessToken
			} else {
				ctx.Error("Error Access Token", fasthttp.StatusForbidden)
			}
		}

		return ""
	}

	ctx.Redirect(f.authCodeURL(f.RequestURL(ctx)), fasthttp.StatusFound)
	return ""
}

func (f *Feishu) retrieveTenantToken() *FeishuTenantToken {
	data := map[string]string{
		"app_id":     f.Config.AppId,
		"app_secret": f.Config.AppSecret,
	}

	body, err := f.postURL(FeishuTenantURL, data, "")
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
	data := map[string]string{
		"grant_type": "authorization_code",
		"code":       code,
	}

	body, err := f.postURL(FeishuTokenURL, data, f.TenantToken.TenantAccessToken)
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

	req.SetRequestURI(url)
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

func (f *Feishu) authCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString(FeishuAuthURL)
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
