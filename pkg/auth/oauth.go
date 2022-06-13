package auth

import "github.com/valyala/fasthttp"

type OAuth struct {
	AppKey         string
	AppSecret      string
	AccessTokenUrl string
	AuthorizeUrl   string
}

func (oauth *OAuth) Handler(ctx *fasthttp.RequestCtx) {

}
