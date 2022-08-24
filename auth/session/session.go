package session

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/joyqi/ngate/config"
	"github.com/valyala/fasthttp"
)

type Session struct {
	cookie *securecookie.SecureCookie
	cfg    *config.SessionConfig
}

func New(cfg config.SessionConfig) (*Session, error) {
	if hashKeyLen := len(cfg.HashKey); hashKeyLen < 32 {
		return nil, fmt.Errorf("hash keys should be at least 32 bytes long")
	}

	if blockKeyLen := len(cfg.BlockKey); blockKeyLen != 16 && blockKeyLen != 32 {
		return nil, fmt.Errorf("block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long")
	}

	return &Session{
		cookie: securecookie.New([]byte(cfg.HashKey), []byte(cfg.BlockKey)),
		cfg:    &cfg,
	}, nil
}

// Store returns a new session store
func (s *Session) Store(ctx *fasthttp.RequestCtx) *Store {
	store := &Store{
		cookie: s.cookie,
		ctx:    ctx,
		cfg:    s.cfg,
	}

	return store
}
