package pipe

import (
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
)

type Frontend struct {
	Addr        string
	BackendAddr string
	Auth        auth.Auth
	Session     *Session
}

func (frontend *Frontend) Serve() {
	log.Success("http pipe %s -> %s", frontend.Addr, frontend.BackendAddr)

	if err := fasthttp.ListenAndServe(frontend.Addr, frontend.handler); err != nil {
		log.Fatal("http server error: %s", err)
	}
}

func (frontend *Frontend) handler(ctx *fasthttp.RequestCtx) {
	session := frontend.Session.Store(ctx)

	defer frontend.close(ctx, session)

	if session.Get("token") != "" {
		frontend.requestBackend(session, ctx)
	} else {
		token, requestUrl := frontend.Auth.Handler(ctx)

		if len(requestUrl) > 0 {
			session.Set("referer", string(ctx.RequestURI()))
		}

		if len(token) > 0 {
			session.Set("token", token)
			frontend.requestBackend(session, ctx)
		}
	}
}

func (frontend *Frontend) requestBackend(session *SessionStore, ctx *fasthttp.RequestCtx) {
	if referer := session.Get("referer"); referer != "" {
		session.Delete("referer")
		ctx.Redirect(referer, fasthttp.StatusFound)
	} else {
		hc := &fasthttp.HostClient{
			Addr: frontend.BackendAddr,
		}

		req := &ctx.Request
		resp := &ctx.Response

		if err := hc.Do(req, resp); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		} else {
			log.Info("%s %s%s %d", req.Header.Method(), req.Host(), req.RequestURI(), resp.StatusCode())
		}
	}
}

func (frontend *Frontend) close(ctx *fasthttp.RequestCtx, session *SessionStore) {
	session.Save()
	log.Info("%s %s%s %d", ctx.Request.Header.Method(), ctx.Request.Host(), ctx.Request.RequestURI(), ctx.Response.StatusCode())
}
