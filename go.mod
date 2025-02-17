module github.com/joyqi/ngate

go 1.21

toolchain go1.23.3

require (
	github.com/gorilla/securecookie v1.1.2
	github.com/valyala/fasthttp v1.58.0
	github.com/yeqown/fasthttp-reverse-proxy/v2 v2.2.3
	golang.org/x/oauth2 v0.26.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/fasthttp/websocket v1.5.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
)

replace github.com/yeqown/fasthttp-reverse-proxy/v2 => ../fasthttp-reverse-proxy
