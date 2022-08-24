package server

import (
	"fmt"
	"github.com/joyqi/ngate/auth"
	"github.com/joyqi/ngate/config"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
	"strings"
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
			GroupValid:      groupValid(pipeConfig.Access),
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

func selectGroups(hostName string, accessConfig []config.AccessConfig) []string {
	for _, c := range accessConfig {
		if c.HostName == hostName {
			return c.Groups
		}
	}

	return nil
}

func existsGroup(group string, selectGroups []string) bool {
	for _, g := range selectGroups {
		if group == g {
			return true
		}
	}

	return false
}

func groupValid(accessConfig []config.AccessConfig) auth.PipeGroupValid {
	return func(group string, hostName string) bool {
		sg := selectGroups(hostName, accessConfig)

		if sg != nil {
			for _, g := range strings.Split(group, ",") {
				if existsGroup(g, sg) {
					return true
				}
			}

			return false
		}

		return true
	}
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
