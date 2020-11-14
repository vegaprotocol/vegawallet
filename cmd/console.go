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
	s          *http.Server
}

func newConsoleProxy(log *zap.Logger, port int, consoleURL string) *consoleProxy {
	return &consoleProxy{
		log:        log,
		port:       port,
		consoleURL: consoleURL,
	}
}

func (c *consoleProxy) Start() error {
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = c.consoleURL
			req.Host = c.consoleURL
		},
	}
	consoleProxyAddr := fmt.Sprintf("localhost:%v", c.port)
	c.s = &http.Server{
		Addr:    consoleProxyAddr,
		Handler: &proxy,
	}

	c.log.Info("starting console proxy",
		zap.String("proxy.address", consoleProxyAddr),
		zap.String("address", c.consoleURL),
	)
	return c.s.ListenAndServe()
}

func (c *consoleProxy) Stop() error {
	return c.s.Shutdown(context.Background())
}

func (c *consoleProxy) GetBrowserURL() string {
	return fmt.Sprintf("http://localhost:%v", c.port)
}
