package pipe

import (
	"fmt"
	"github.com/joyqi/ngate/internal/config"
	"github.com/joyqi/ngate/internal/log"
	"github.com/joyqi/ngate/pkg/auth"
	proxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
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

		httpProxy, proxyErr := proxy.NewReverseProxyWith(
			proxy.WithAddress(backendAddr),
			proxy.WithDisablePathNormalizing(true),
			proxy.WithDisableVirtualHost(true),
			proxy.WithDebug(),
			proxy.WithTimeout(time.Duration(pipeConfig.Backend.Timeout)*time.Millisecond))

		if proxyErr != nil {
			return proxyErr
		}

		wsProxies := sync.Map{}

		getter := func(path string) (*proxy.WSReverseProxy, error) {
			if v, ok := wsProxies.Load(path); !ok {
				wsProxy, err := proxy.NewWSReverseProxyWith(
					proxy.WithURL_OptionWS("ws://" + backendAddr + path))

				if err != nil {
					return nil, err
				}

				wsProxies.Store(path, wsProxy)
				return wsProxy, nil
			} else {
				return v.(*proxy.WSReverseProxy), nil
			}
		}

		frontend := &Frontend{
			Addr:                 addr,
			Session:              session,
			Auth:                 auth,
			Wait:                 wg,
			GroupValid:           groupValid(pipeConfig.Access),
			BackendHostName:      pipeConfig.Backend.HostName,
			WSBackendProxies:     &wsProxies,
			WSBackendProxyGetter: getter,
			BackendProxy:         httpProxy,
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
