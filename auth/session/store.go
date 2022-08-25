package session

import (
	"encoding/json"
	"github.com/gorilla/securecookie"
	"github.com/joyqi/ngate/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Store struct {
	cookie *securecookie.SecureCookie
	ctx    *fasthttp.RequestCtx
	cfg    *config.SessionConfig
	stored bool

	// Token stores the token
	Token *Token
}

// Init session store from cookie
func (store *Store) Init() {
	cookie := store.ctx.Request.Header.Cookie(store.cfg.CookieKey)
	data := make(map[string]string)

	if len(cookie) > 0 {
		if err := store.cookie.Decode(store.cfg.CookieKey, string(cookie), &data); err != nil {
			log.Warning("wrong cookie: %s", err)
			return
		} else {
			log.Debug("cookie: %s", data)
		}

		if value, ok := data["_"]; !ok {
			log.Warning("no _ in cookie")
			return
		} else {
			token := &Token{}
			if err := json.Unmarshal([]byte(value), token); err != nil {
				log.Warning("wrong token: %s", err)
			}
		}
	}
}

// Commit changes store flag
func (store *Store) Commit() {
	store.stored = true
}

func (store *Store) Save() {
	if !store.stored {
		return
	}

	var c fasthttp.Cookie

	if store.Token == nil {
		c.SetExpire(fasthttp.CookieExpireDelete)
	} else {
		data := make(map[string]string)
		value, err := json.Marshal(store.Token)

		if err != nil {
			log.Warning("json encode error")
			return
		}

		data["_"] = string(value)
		encoded, err := store.cookie.Encode(store.cfg.CookieKey, &data)

		if err != nil {
			log.Warning("encrypt error %s", err)
			return
		}

		c.SetKey(store.cfg.CookieKey)
		c.SetValue(encoded)
		c.SetPath("/")
		c.SetDomain(store.cfg.CookieDomain)
		c.SetHTTPOnly(true)
		c.SetSecure(true)
		c.SetSameSite(fasthttp.CookieSameSiteNoneMode)
		store.ctx.Response.Header.SetCookie(&c)
	}
}
