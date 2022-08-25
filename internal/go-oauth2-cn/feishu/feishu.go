package feishu

import (
	"bytes"
	"context"
	"net/url"
	"sync"
	"time"
)

var EndpointURL = Endpoint{
	AuthURL:          "https://open.feishu.cn/open-apis/authen/v1/index",
	TokenURL:         "https://open.feishu.cn/open-apis/authen/v1/access_token",
	RefreshTokenURL:  "https://open.feishu.cn/open-apis/authen/v1/refresh_access_token",
	UserGroupsApiURL: "https://open.feishu.cn/open-apis/contact/v3/group/member_belong",
}

// Endpoint is the endpoint to connect to the server.
type Endpoint struct {
	AuthURL  string
	TokenURL string

	// RefreshTokenURL is the URL to refresh the token.
	RefreshTokenURL string

	// UserGroupsApiURL is the URL to retrieve user groups.
	UserGroupsApiURL string
}

// Config represents the configuration of the feishu service
type Config struct {
	// AppID is the app id of feishu.
	AppID string

	// AppSecret is the app secret of feishu.
	AppSecret string

	// RedirectURL is the URL to redirect users going through
	RedirectURL string

	// tenantTokenMu is the lock for tenant token request
	tenantTokenMu sync.Mutex

	// tenantToken is the tenant token
	tenantToken string

	// tenantTokenExpireAt is the expiration time of the tenant token
	tenantTokenExpireAt time.Time
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
func (c *Config) Exchange(ctx context.Context, code string) (*Token, error) {
	req := &TokenRequest{
		GrantType: "authorization_code",
		Code:      code,
	}

	return retrieveToken(EndpointURL.TokenURL, req, c)
}

// TokenSource returns a TokenSource to grant token access
func (c *Config) TokenSource(ctx context.Context, t *Token) *TokenSource {
	return &TokenSource{
		ctx:  ctx,
		conf: c,
		t:    t,
	}
}