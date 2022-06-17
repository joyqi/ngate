package pipe

import (
	"fmt"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/pkg/auth"
)

type Pipe interface {
	Serve(auth auth.Auth, bc config.BackendConfig)
}

func New(cfg *config.Config, auth auth.Auth) error {
	if len(cfg.Pipes) == 0 {
		return fmt.Errorf("empty pipes")
	}

	for _, pipeConfig := range cfg.Pipes {
		session, err := NewSession(pipeConfig.Session)
		if err != nil {
			return err
		}

		frontend := &Frontend{
			Addr:           fmt.Sprint(Addr{pipeConfig.Host, pipeConfig.Port}),
			BackendAddr:    fmt.Sprint(Addr{pipeConfig.Backend.Host, pipeConfig.Backend.Port}),
			BackendTimeout: pipeConfig.Backend.Timeout,
			Session:        session,
			Auth:           auth,
		}

		err = frontend.Serve()
		if err != nil {
			return err
		}
	}

	return nil
}

type Addr struct {
	Host string
	Port int
}

func (a Addr) String() string {
	if a.Host == "" {
		a.Host = "0.0.0.0"
	}

	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
