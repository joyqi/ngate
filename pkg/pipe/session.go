package pipe

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type Session struct {
	Cookie *securecookie.SecureCookie
	Config *config.SessionConfig
}

func NewSession(cfg config.SessionConfig) *Session {
	if hashKeyLen := len(cfg.HashKey); hashKeyLen < 32 {
		log.Fatal("hash keys should be at least 32 bytes long")
	}

	if blockKeyLen := len(cfg.BlockKey); blockKeyLen != 16 && blockKeyLen != 32 {
		log.Fatal("block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long")
	}

	return &Session{
		Cookie: securecookie.New([]byte(cfg.HashKey), []byte(cfg.BlockKey)),
		Config: &cfg,
	}
}

func (session *Session) Store(ctx *fasthttp.RequestCtx) *SessionStore {
	store := &SessionStore{
		Cookie: session.Cookie,
		Ctx:    ctx,
		Config: session.Config,
		Stored: false,
	}

	store.Init()
	return store
}

type SessionStore struct {
	Cookie *securecookie.SecureCookie
	Ctx    *fasthttp.RequestCtx
	Config *config.SessionConfig
	Stored bool
	data   map[string]string
}

// Init session store from cookie
func (store *SessionStore) Init() {
	cookie := store.Ctx.Request.Header.Cookie(store.Config.CookieKey)
	store.data = make(map[string]string)

	if len(cookie) > 0 {
		if err := store.Cookie.Decode(store.Config.CookieKey, string(cookie), &store.data); err != nil {
			log.Warning("wrong cookie: %s", err)
		}
	}
}

func (store *SessionStore) Set(key string, value string) {
	store.data[key] = value
	store.Stored = true
}

func (store *SessionStore) SetInt(key string, value int64) {
	store.Set(key, fmt.Sprintf("%d", value))
}

// Get string value from session store
func (store *SessionStore) Get(key string) string {
	value, ok := store.data[key]

	if !ok {
		return ""
	}

	return value
}

// GetInt int value from session store
func (store *SessionStore) GetInt(key string) int64 {
	value := store.Get(key)

	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}

	return 0
}

func (store *SessionStore) Expired(last int64) bool {
	return time.Now().Unix() > last+store.Config.ExpiresIn
}

func (store *SessionStore) Delete(key string) {
	delete(store.data, key)
	store.Stored = true
}

func (store *SessionStore) Save() {
	if !store.Stored {
		return
	}

	var c fasthttp.Cookie
	value, err := store.Cookie.Encode(store.Config.CookieKey, &store.data)

	if err == nil {
		c.SetKey(store.Config.CookieKey)
		c.SetValue(value)
		c.SetPath("/")
		c.SetDomain(store.Config.CookieDomain)
		store.Ctx.Response.Header.SetCookie(&c)
	}
}
