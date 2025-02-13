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
	GroupValid           auth.PipeGroupValid
	Addr                 string
	Auth                 auth.Auth
	Session              *Session
	Wait                 *sync.WaitGroup
	BackendHostName      string
	WSBackendProxies     *sync.Map
	WSBackendProxyGetter func(path string) (*proxy.WSReverseProxy, error)
	BackendProxy         *proxy.ReverseProxy
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopHeaders = []string{
	"Connection",          // Connection
	"Proxy-Connection",    // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",          // Keep-Alive
	"Proxy-Authenticate",  // Proxy-Authenticate
	"Proxy-Authorization", // Proxy-Authorization
	"Te",                  // canonicalized version of "TE"
	"Trailer",             // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",   // Transfer-Encoding
	"Upgrade",             // Upgrade
}

func (frontend *Frontend) Serve() {
	defer frontend.Wait.Done()

	s := fasthttp.Server{
		Handler:                       frontend.handler,
		ReadBufferSize:                64 * 1024,
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

	if frontend.Auth.Valid(session) {
		if frontend.Auth.GroupValid(string(ctx.Request.Host()), session, frontend.GroupValid) {
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
		wsProxy, err := frontend.WSBackendProxyGetter(path)

		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		wsProxy.ServeHTTP(ctx)
	} else {
		frontend.BackendProxy.ServeHTTP(ctx)
	}

	// for _, h := range hopHeaders {
	//	req.Header.Del(h)
	// }

	/*
		if err := frontend.BackendProxy.DoTimeout(req, resp, frontend.BackendTimeout); err != nil {
			if errors.Is(err, fasthttp.ErrTimeout) {
				ctx.Error(err.Error(), fasthttp.StatusRequestTimeout)
			} else {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			}
			log.Error("%s %s%s %s", req.Header.Method(), req.Host(), req.RequestURI(), err.Error())
		}
	*/

	// for _, h := range hopHeaders {
	//	resp.Header.Del(h)
	// }
}

func (frontend *Frontend) close(ctx *fasthttp.RequestCtx, session *SessionStore) {
	session.Save()
	log.Info("%s %s%s %d", ctx.Request.Header.Method(), ctx.Request.Host(), ctx.Request.RequestURI(), ctx.Response.StatusCode())
}
