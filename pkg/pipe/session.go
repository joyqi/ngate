package pipe

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
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
		Valid:  false,
	}

	store.Init()
	return store
}

type SessionStore struct {
	Cookie *securecookie.SecureCookie
	Ctx    *fasthttp.RequestCtx
	Config *config.SessionConfig
	Valid  bool
	data   map[string]string
}

// Init session store from cookie
func (store *SessionStore) Init() {
	cookie := store.Ctx.Request.Header.Cookie(store.Config.CookieKey)
	store.data = make(map[string]string)

	if err := store.Cookie.Decode(store.Config.CookieKey, string(cookie), &store.data); err == nil {
		now := time.Now().Unix()
		lastTime := store.GetInt("time")

		if lastTime == 0 || int64(lastTime+store.Config.ExpireHours*3600) >= now {
			store.Valid = true
		}
	}
}

func (store *SessionStore) Set(key string, value string) {
	store.data[key] = value
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
func (store *SessionStore) GetInt(key string) int {
	value := store.Get(key)

	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	return 0
}

func (store *SessionStore) Delete(key string) {
	delete(store.data, key)
}

func (store *SessionStore) Save() {
	store.SetInt("time", time.Now().Unix()+int64(store.Config.ExpireHours*3600))

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
