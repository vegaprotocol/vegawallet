package cmd

import (
	"fmt"
	"net/url"
	"os"

	"go.uber.org/zap"
)

func init() {
	err := zap.RegisterSink("winfile", newWinFileSink)
	if err != nil {
		panic(fmt.Errorf("couldn't register the windows file sink: %w", err))
	}
}

func toZapLogPath(logPath string) string {
	return "winfile:///" + logPath
}

// newWinFileSink creates a log sink on Windows machines as zap, by default,
// doesn't support Windows paths. A workaround is to create a fake winfile
// scheme and register it with zap instead. This workaround is taken from
// the GitHub issue at https://github.com/uber-go/zap/issues/621.
func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse().
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644) // nolint:gomnd
}
