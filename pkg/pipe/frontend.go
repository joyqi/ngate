package pipe

import (
	"encoding/json"
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

	if session.Get("token") != "" {
		frontend.requestBackend(ctx)
	} else {
		token := frontend.Auth.Handler(ctx, frontend.SoftRedirect(ctx))

		if len(token) > 0 {
			session.Set("token", token)
		}
	}
}

func (frontend *Frontend) requestBackend(ctx *fasthttp.RequestCtx) {
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

func (frontend *Frontend) close(ctx *fasthttp.RequestCtx, session *SessionStore) {
	session.Save()
	log.Info("%s %s%s %d", ctx.Request.Header.Method(), ctx.Request.Host(), ctx.Request.RequestURI(), ctx.Response.StatusCode())
}
