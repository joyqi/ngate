package pipe

import (
	"github.com/joyqi/dahuang/pkg/auth"
	"github.com/joyqi/dahuang/pkg/config"
	"github.com/joyqi/dahuang/pkg/log"
)

type Pipe interface {
	Serve(auth auth.Auth)
}

func New(cfg *config.Config, auth auth.Auth) {
	if len(cfg.Pipes) == 0 {
		log.Fatal("empty pipes")
	}

	for _, pipeConfig := range cfg.Pipes {
		var pipe Pipe

		switch pipeConfig.Kind {
		case "http":
		default:
			pipe = &Http{
				Host: pipeConfig.Host,
				Port: pipeConfig.Port,
			}
		}

		pipe.Serve(auth)
	}
}
