package pipe

import (
	"fmt"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/pkg/auth"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"sync"
	"time"
)

type Pipe interface {
	Serve(auth auth.Auth, bc config.BackendConfig)
}

func New(cfg *config.Config, auth auth.Auth) error {
	if len(cfg.Pipes) == 0 {
		return fmt.Errorf("empty pipes")
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.Pipes))

	for _, pipeConfig := range cfg.Pipes {
		session, err := NewSession(pipeConfig.Session)
		if err != nil {
			return err
		}

		addr := net.JoinHostPort(
			defaultHost(pipeConfig.Host, "0.0.0.0"),
			defaultPort(pipeConfig.Port, 8000))

		backendAddr := net.JoinHostPort(
			defaultHost(pipeConfig.Backend.Host, "127.0.0.1"),
			defaultPort(pipeConfig.Backend.Port, 8000))

		frontend := &Frontend{
			Addr:            addr,
			Session:         session,
			Auth:            auth,
			Wait:            wg,
			BackendHostName: pipeConfig.Backend.HostName,
			BackendTimeout:  time.Duration(pipeConfig.Backend.Timeout) * time.Millisecond,
			BackendProxy: &fasthttp.HostClient{
				Addr:                          backendAddr,
				DisableHeaderNamesNormalizing: true,
				DisablePathNormalizing:        true,
				ReadBufferSize:                64 * 1024,
			},
		}

		go frontend.Serve()
	}

	wg.Wait()
	return nil
}

func defaultHost(host string, defaultHost string) string {
	if host == "" {
		return defaultHost
	}

	return host
}

func defaultPort(port int, defaultPort int) string {
	if port == 0 {
		port = defaultPort
	}

	return strconv.Itoa(port)
}
