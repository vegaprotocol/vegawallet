//go:build !windows
// +build !windows

package cmd

func toZapLogPath(p string) string {
	return p
}
