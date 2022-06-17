package pipe

import (
	"fmt"
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
)

type Pipe interface {
	Serve(auth auth.Auth, bc config.BackendConfig)
}

func New(cfg *config.Config, auth auth.Auth) {
	if len(cfg.Pipes) == 0 {
		log.Fatal("empty pipes")
	}

	for _, pipeConfig := range cfg.Pipes {
		frontend := &Frontend{
			Addr:           fmt.Sprint(Addr{pipeConfig.Host, pipeConfig.Port}),
			BackendAddr:    fmt.Sprint(Addr{pipeConfig.Backend.Host, pipeConfig.Backend.Port}),
			BackendTimeout: pipeConfig.Backend.Timeout,
			Session:        NewSession(pipeConfig.Session),
			Auth:           auth,
		}

		frontend.Serve()
	}
}

type Addr struct {
	Host string
	Port int
}

func (a Addr) String() string {
	if a.Host == "" {
		a.Host = "0.0.0.0"
	}

	if a.Port <= 0 {
		log.Fatal("wrong port number %d for %s", a.Port, a.Host)
	}

	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
