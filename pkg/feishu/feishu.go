package feishu

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/joyqi/ngate/pkg/http"
	"github.com/valyala/fasthttp"
	"net/url"
)

var EndpointURL = Endpoint{
	AuthURL:         "https://open.feishu.cn/open-apis/authen/v1/index",
	TokenURL:        "https://open.feishu.cn/open-apis/authen/v1/access_token",
	RefreshTokenURL: "https://open.feishu.cn/open-apis/authen/v1/refresh_access_token",
	UserGroupApiURL: "https://open.feishu.cn/open-apis/contact/v3/group/member_belong",
}

// Endpoint is the endpoint to connect to the server.
type Endpoint struct {
	AuthURL  string
	TokenURL string

	// RefreshTokenURL is the URL to refresh the token.
	RefreshTokenURL string

	// TenantTokenURL is the URL to request the tenant token.
	TenantTokenURL string

	// USerGroupApiURL is the URL to request the user group API.
	UserGroupApiURL string
}

// Token represents the credentials
type Token struct {
	// AccessToken is the token used to access the application
	AccessToken string `json:"access_token"`

	// RefreshToken is the token used to refresh the user's access token
	RefreshToken string `json:"refresh_token"`

	// ExpiresIn is the number of seconds the token will be valid
	ExpiresIn int64 `json:"expires_in"`
}

type TokenRequest struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

// TokenResponse represents the response from the Token service
type TokenResponse struct {
	// Code is the response status code
	Code int `json:"code"`

	// Msg is the response message in the response body
	Msg string `json:"msg"`

	// Data is the response body data
	Data Token `json:"data"`
}

// Config represents the configuration of the feishu service
type Config struct {
	// AppID is the app id of feishu.
	AppID string

	// AppSecret is the app secret of feishu.
	AppSecret string

	// RedirectURL is the URL to redirect users going through
	RedirectURL string
}

// AuthCodeURL is the URL to redirect users going through authentication
func (c *Config) AuthCodeURL(state string) string {
	var buf bytes.Buffer
	buf.WriteString(EndpointURL.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"app_id":        {c.AppID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}

	if state != "" {
		v.Set("state", state)
	}

	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

// Exchange retrieve the token from access token endpoint
func (c *Config) Exchange(ctx *fasthttp.RequestCtx, code string, tenantToken string) (*Token, error) {
	req := &TokenRequest{
		GrantType: "authorization_code",
		Code:      code,
	}

	body, err := http.PostJSON(
		EndpointURL.TokenURL,
		req,
		http.Header{Key: "Authorization", Value: tenantToken},
	)

	if err != nil {
		return nil, err
	}

	resp := TokenResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Msg)
	}

	return &resp.Data, nil
}
