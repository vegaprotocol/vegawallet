package console

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"code.vegaprotocol.io/go-wallet/version"
)

type Console struct {
	server      *http.Server
}

func NewConsole(port int, consoleURL, nodeURL string) *Console {
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Set("Referer", nodeURL)
			req.Header.Set(
				"User-Agent",
				fmt.Sprintf("%v VegaWallet/%v", req.Header.Get("User-Agent"), version.Version),
			)
			req.Header.Set("X-Vega-Wallet-Version", version.Version)

			// To prevent IP spoofing, be sure to delete any pre-existing
			// X-Forwarded-For header coming from the client
			req.Header.Set("X-Forwarded-For", "")

			req.URL.Scheme = "https"
			req.URL.Host = consoleURL
			req.Host = consoleURL
		},
	}

	addr := fmt.Sprintf("127.0.0.1:%v", port)
	return &Console{
		server: &http.Server{
			Addr:    addr,
			Handler: &proxy,
		},
	}
}

func (c *Console) Start() error {
	return c.server.ListenAndServe()
}

func (c *Console) Stop() error {
	return c.server.Shutdown(context.Background())
}

func (c *Console) GetBrowserURL() string {
	return fmt.Sprintf("http://%s", c.server.Addr)
}
