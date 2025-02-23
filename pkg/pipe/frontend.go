package pipe

import (
	"encoding/json"
	"github.com/joyqi/ngate/internal/log"
	"github.com/joyqi/ngate/pkg/auth"
	"github.com/valyala/fasthttp"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
	"sync"
)

type Frontend struct {
	GroupValid      auth.PipeGroupValid
	Addr            string
	Auth            auth.Auth
	Session         *Session
	Wait            *sync.WaitGroup
	BackendHostName string
	WSBackendProxy  *proxy.WSReverseProxy
	BackendProxy    *proxy.ReverseProxy
}

func (frontend *Frontend) Serve() {
	defer frontend.Wait.Done()

	s := fasthttp.Server{
		Handler:                       frontend.handler,
		ReadBufferSize:                64 * 1024,
		StreamRequestBody:             true,
		DisablePreParseMultipartForm:  true,
		DisableHeaderNamesNormalizing: true,
	}

	if err := s.ListenAndServe(frontend.Addr); err != nil {
		log.Error("http server error: %s", err)
	}
}

// SoftRedirect perform a redirect handle by javascript code
func (frontend *Frontend) SoftRedirect(ctx *fasthttp.RequestCtx) auth.SoftRedirect {
	return func(url string) {
		u, err := json.Marshal(url)
		if err != nil {
			ctx.Error("Wrong Url", fasthttp.StatusBadRequest)
			return
		}

		ctx.SetContentType("text/html")
		if _, err = ctx.WriteString("<script>window.location.href=" + string(u) + "</script>"); err != nil {
			ctx.Error("Wrong Url", fasthttp.StatusBadRequest)
		}
	}
}

func (frontend *Frontend) handler(ctx *fasthttp.RequestCtx) {
	session := frontend.Session.Store(ctx)

	defer frontend.close(ctx, session)

	if frontend.Auth.Valid(ctx, session) {
		if frontend.GroupValid(frontend.Auth.Groups(ctx, session), string(ctx.Request.Host())) {
			frontend.requestBackend(ctx)
		} else {
			ctx.Error("Access Denied", fasthttp.StatusForbidden)
		}
	} else {
		frontend.Auth.Handler(ctx, session, frontend.SoftRedirect(ctx))
	}
}

func (frontend *Frontend) requestBackend(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	path := string(ctx.Path())

	// set hostname for request
	if frontend.BackendHostName != "" {
		req.SetHost(frontend.BackendHostName)
	}

	// detect if the request is websocket
	if string(req.Header.Peek("Upgrade")) == "websocket" {
		req.Header.Set(proxy.DefaultOverrideHeader, path)
		frontend.WSBackendProxy.ServeHTTP(ctx)
	} else {
		frontend.BackendProxy.ServeHTTP(ctx)
	}
}

func (frontend *Frontend) close(ctx *fasthttp.RequestCtx, session *SessionStore) {
	session.Save()
	log.Info("%s %s%s %d", ctx.Request.Header.Method(), ctx.Request.Host(), ctx.Request.RequestURI(), ctx.Response.StatusCode())
}
