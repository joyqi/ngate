package pipe

import "github.com/joyqi/dahuang/pkg/log"

type Http struct {
	Host string
	Port int
}

func (h *Http) Serve() {
	var host string

	if h.Host == "" {
		host = "0.0.0.0"
	} else {
		host = h.Host
	}

	if h.Port <= 0 {
		log.Fatal("wrong port number %d for %s", h.Port, host)
	}
}
