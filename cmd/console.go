package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"go.uber.org/zap"
)

type consoleProxy struct {
	log        *zap.Logger
	port       int
	consoleURL string
	nodeURL    string
	s          *http.Server
	version    string
}

func newConsoleProxy(log *zap.Logger, port int, consoleURL, nodeURL, version string) *consoleProxy {
	return &consoleProxy{
		log:        log,
		port:       port,
		consoleURL: consoleURL,
		nodeURL:    nodeURL,
		version:    version,
	}
}

func (c *consoleProxy) Start() error {
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("Referer", c.nodeURL)
			req.Header.Set(
				"User-Agent",
				fmt.Sprintf("%v VegaWallet/%v", req.Header.Get("User-Agent"), c.version),
			)
			req.URL.Scheme = "https"
			req.URL.Host = c.consoleURL
			req.Host = c.consoleURL
		},
	}
	consoleProxyAddr := fmt.Sprintf("127.0.0.1:%v", c.port)
	c.s = &http.Server{
		Addr:    consoleProxyAddr,
		Handler: &proxy,
	}

	// c.log.Info("starting console proxy",
	// 	zap.String("proxy.address", consoleProxyAddr),
	// 	zap.String("address", c.consoleURL),
	// )
	return c.s.ListenAndServe()
}

func (c *consoleProxy) Stop() error {
	return c.s.Shutdown(context.Background())
}

func (c *consoleProxy) GetBrowserURL() string {
	return fmt.Sprintf("http://127.0.0.1:%v", c.port)
}
