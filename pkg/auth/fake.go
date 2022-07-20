package auth

import (
	"github.com/valyala/fasthttp"
)

type Fake struct {
}

func (f *Fake) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	// do nothing
}

func (f *Fake) Valid(session Session) bool {
	return true
}

func (f *Fake) GroupValid(hostName string, session Session, valid PipeGroupValid) bool {
	return true
}
