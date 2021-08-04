package console

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"code.vegaprotocol.io/go-wallet/version"
)

type Console struct {
	port       int
	consoleURL string
	nodeURL    string
	server     *http.Server
}

func NewConsole(port int, consoleURL, nodeURL string) *Console {
	return &Console{
		port:       port,
		consoleURL: consoleURL,
		nodeURL:    nodeURL,
	}
}

func (c *Console) Start() error {
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("Referer", c.nodeURL)
			req.Header.Set(
				"User-Agent",
				fmt.Sprintf("%v VegaWallet/%v", req.Header.Get("User-Agent"), version.Version),
			)
			req.Header.Set("X-Vega-Wallet-Version", version.Version)

			// To prevent IP spoofing, be sure to delete any pre-existing
			// X-Forwarded-For header coming from the client
			req.Header.Set("X-Forwarded-For", "")

			req.URL.Scheme = "https"
			req.URL.Host = c.consoleURL
			req.Host = c.consoleURL
		},
	}
	consoleProxyAddr := fmt.Sprintf("127.0.0.1:%v", c.port)
	c.server = &http.Server{
		Addr:    consoleProxyAddr,
		Handler: &proxy,
	}

	return c.server.ListenAndServe()
}

func (c *Console) Stop() error {
	return c.server.Shutdown(context.Background())
}

func (c *Console) GetBrowserURL() string {
	return c.server.Addr
}
