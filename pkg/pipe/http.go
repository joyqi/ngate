package pipe

import (
	"fmt"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/log"
	"github.com/valyala/fasthttp"
)

type Http struct {
	Host string
	Port int
}

func (h *Http) Serve(auth auth.Auth) {
	var host string

	if h.Host == "" {
		host = "0.0.0.0"
	} else {
		host = h.Host
	}

	if h.Port <= 0 {
		log.Fatal("wrong port number %d for %s", h.Port, host)
	}

	handler := Handler{auth}
	addr := fmt.Sprintf("%s:%d", h.Host, h.Port)
	log.Info("http server listening at %s", addr)

	err := fasthttp.ListenAndServe(addr, handler.Handler)
	if err != nil {
		log.Fatal("http server error: %s", err)
	}
}

type Handler struct {
	auth.Auth
}

func (h *Handler) Handler(ctx *fasthttp.RequestCtx) {

}
