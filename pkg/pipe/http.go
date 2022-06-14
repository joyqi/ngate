package pipe

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

const (
	CookieKey         = "__dh_auth"
	CookieRedirectKey = CookieKey + "_redirect"
)

type Http struct {
	Addr   string
	Cookie HttpCookie
}

type HttpCookie struct {
	HashKey       []byte
	BlockKey      []byte
	ExpireSeconds int
}

func (h *Http) Serve(auth auth.Auth, bc config.BackendConfig) {
	backend := HttpBackend{
		auth,
		securecookie.New(h.Cookie.HashKey, h.Cookie.BlockKey),
		h.Cookie.ExpireSeconds,
		fmt.Sprint(Addr{bc.Host, bc.Port}),
	}

	h.validate(backend.Addr)
	log.Success("http pipe %s -> %s", h.Addr, backend.Addr)

	if err := fasthttp.ListenAndServe(h.Addr, backend.Serve); err != nil {
		log.Fatal("http server error: %s", err)
	}
}

func (h *Http) validate(backendAddr string) {
	if h.Addr == backendAddr {
		log.Fatal("the server address can not be equal to the backend address")
	}

	if hashKeyLen := len(h.Cookie.HashKey); hashKeyLen < 32 {
		log.Fatal("hash keys should be at least 32 bytes long")
	}

	if blockKeyLen := len(h.Cookie.BlockKey); blockKeyLen != 16 && blockKeyLen != 32 {
		log.Fatal("block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long")
	}
}

type HttpBackend struct {
	auth.Auth
	Cookie        *securecookie.SecureCookie
	ExpireSeconds int
	Addr          string
}

// Serve handle all http requests from clients
func (h *HttpBackend) Serve(ctx *fasthttp.RequestCtx) {
	if h.verifyCookie(ctx) {
		if !h.recoverUrl(ctx) {
			h.proxyPass(ctx)
		}
	} else {
		if !h.Auth.Handler(ctx) {
			h.rememberUrl(ctx)
		}
	}
}

func (h *HttpBackend) verifyCookie(ctx *fasthttp.RequestCtx) bool {
	c := ctx.Request.Header.Cookie(CookieKey)
	if len(c) == 0 {
		return false
	}

	cookie := make(map[string]string)
	if err := h.Cookie.Decode(CookieKey, string(c), &cookie); err == nil {
		value, ok := cookie["time"]

		if ok {
			if i, err := strconv.Atoi(value); err == nil {
				now := time.Now().Unix()

				if int64(i+h.ExpireSeconds) >= now {
					return true
				}
			}
		}
	}

	return false
}

func (h *HttpBackend) rememberUrl(ctx *fasthttp.RequestCtx) {
	var c fasthttp.Cookie

	c.SetKey(CookieRedirectKey)
	c.SetValueBytes(ctx.Request.RequestURI())
	ctx.Response.Header.Cookie(&c)
}

func (h *HttpBackend) recoverUrl(ctx *fasthttp.RequestCtx) bool {
	if c := ctx.Request.Header.Cookie(CookieRedirectKey); len(c) > 0 {
		ctx.Response.Header.DelCookie(CookieRedirectKey)
		ctx.Redirect(string(c), fasthttp.StatusFound)
		return true
	}

	return false
}

func (h *HttpBackend) proxyPass(ctx *fasthttp.RequestCtx) {
	hc := &fasthttp.HostClient{
		Addr: h.Addr,
	}

	req := &ctx.Request
	resp := &ctx.Response

	if err := hc.Do(req, resp); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	} else {
		log.Info("%s %s %d", req.Header.Method(), req.RequestURI(), resp.StatusCode())
	}

	return
}
