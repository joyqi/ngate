module github.com/joyqi/ngate

go 1.21

toolchain go1.23.3

require (
	github.com/gorilla/securecookie v1.1.2
	github.com/joyqi/go-lafi v0.2.0
	github.com/valyala/fasthttp v1.58.0
	github.com/wenerme/go-wecom v0.10.1
	github.com/yeqown/fasthttp-reverse-proxy/v2 v2.2.3
	golang.org/x/oauth2 v0.26.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/creasty/defaults v1.6.0 // indirect
	github.com/fasthttp/websocket v1.5.12 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/wenerme/go-req v0.0.0-20210907160348-d822e81276bb // indirect
	golang.org/x/net v0.35.0 // indirect
)

replace github.com/yeqown/fasthttp-reverse-proxy/v2 => github.com/joyqi/fasthttp-reverse-proxy/v2 v2.0.0-20250213105017-13c262a9116b
