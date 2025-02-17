package auth

import (
	"github.com/valyala/fasthttp"
)

type Fake struct {
}

func NewFake() *Fake {
	return &Fake{}
}

func (f *Fake) Handler(ctx *fasthttp.RequestCtx, session Session, redirect SoftRedirect) {
	// do nothing
}

func (f *Fake) Valid(ctx *fasthttp.RequestCtx, session Session) bool {
	return true
}

func (f *Fake) Groups(ctx *fasthttp.RequestCtx, session Session) []string {
	return nil
}
