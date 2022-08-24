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

	data := make(map[string]string)
}
