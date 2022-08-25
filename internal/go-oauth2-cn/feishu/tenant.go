package feishu

import (
	"encoding/json"
	"errors"
	"github.com/joyqi/ngate/internal/go-oauth2-cn/http"
	"time"
)

var TenantEndpointURL = TenantEndpoint{
	TokenURL: "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
}

type TenantEndpoint struct {
	TokenURL string
}

// TenantTokenRequest represents a request to retrieve a tenant token
type TenantTokenRequest struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// TenantTokenResponse is the token response of the tenant endpoint
type TenantTokenResponse struct {
	// Code is the response status code
	Code int `json:"code"`

	// Msg is the response message
	Msg string `json:"msg"`

	// TenantAccessToken is the access token
	TenantAccessToken string `json:"tenant_access_token"`

	// Expire is the expiration time of the access token
	Expire int64 `json:"expire"`
}

// TenantToken represents a tenant access token from tenant token endpoint
func (c *Config) TenantToken() (string, error) {
	defer c.tenantTokenMu.Unlock()

	c.tenantTokenMu.Lock()
	if !c.TenantTokenValid() {
		req := &TenantTokenRequest{
			AppID:     c.AppID,
			AppSecret: c.AppSecret,
		}

		body, err := http.PostJSON(TenantEndpointURL.TokenURL, req)
		if err != nil {
			return "", err
		}

		resp := TenantTokenResponse{}
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return "", err
		}

		if resp.Code != 0 {
			return "", errors.New(resp.Msg)
		}

		c.tenantToken = resp.TenantAccessToken
		c.tenantTokenExpireAt = time.Now().Add(time.Duration(resp.Expire) * time.Second)
	}

	return c.tenantToken, nil
}

// TenantTokenValid returns true if the tenant token is valid
func (c *Config) TenantTokenValid() bool {
	return c.tenantToken != "" && time.Now().Add(time.Minute).Before(c.tenantTokenExpireAt)
}