package pipe

import (
	"fmt"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/joyqi/ngate/pkg/auth"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
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

		httpProxy, proxyErr := proxy.NewReverseProxyWith(
			proxy.WithAddress(backendAddr),
			proxy.WithDisablePathNormalizing(true),
			proxy.WithDisableVirtualHost(true),
			proxy.WithStreamResponseBody(64*1024),
			proxy.WithTimeout(time.Duration(pipeConfig.Backend.Timeout)*time.Millisecond))

		if proxyErr != nil {
			return proxyErr
		}

		wsProxy, wsProxyErr := proxy.NewWSReverseProxyWith(
			proxy.WithURL_OptionWS("ws://"+backendAddr),
			proxy.WithDynamicPath_OptionWS(true, proxy.DefaultOverrideHeader))

		if wsProxyErr != nil {
			return err
		}

		frontend := &Frontend{
			Addr:            addr,
			Session:         session,
			Auth:            auth,
			Wait:            wg,
			GroupValid:      groupValid(pipeConfig.Access),
			BackendHostName: pipeConfig.Backend.HostName,
			WSBackendProxy:  wsProxy,
			BackendProxy:    httpProxy,
		}

		log.Success("http pipe %s -> %s", addr, backendAddr)
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
	return func(groups []string, hostName string) bool {
		sg := selectGroups(hostName, accessConfig)

		if sg != nil {
			for _, g := range groups {
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
