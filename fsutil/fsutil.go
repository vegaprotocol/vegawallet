package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirPerms = 0700
)

// PathNotFound represent an error when a fd cannot be found
// on the user filesystem
type PathNotFound struct {
	path string
}

// Error return an human readable formating of the error
func (err *PathNotFound) Error() string {
	return fmt.Sprintf("not found: %s", err.path)
}

// DefaultVegaDir returns the location to vega config files and data files:
// binary is in /usr/bin/ -> look for /etc/vega/config.toml
// binary is in /usr/local/vega/bin/ -> look for /usr/local/vega/etc/config.toml
// binary is in /usr/local/bin/ -> look for /usr/local/etc/vega/config.toml
// otherwise, look for $HOME/.vega/config.toml
func DefaultVegaDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	if exPath == "/usr/bin" {
		return "/etc/vega"
	}
	if exPath == "/usr/local/vega/bin" {
		return "/usr/local/vega/etc"
	}
	if exPath == "/usr/local/bin" {
		return "/usr/local/etc/vega"
	}
	return os.ExpandEnv("$HOME/.vega")
}

// EnsureDir will make sure a directory exists or is created at a given filesystem path.
func EnsureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, dirPerms)
		}
		return err
	}
	return nil
}

// PathExists returns whether a link exists at a given filesystem path.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, &PathNotFound{path}
	}
	return false, err
}

// FileExists similar to PathExists, but ensures the path is to a file, not a directory
func FileExists(path string) (bool, error) {
	fs, err := os.Stat(path)
	if err == nil {
		// fileStat -> is not a directory
		ok := !fs.IsDir()
		return ok, nil
	}
	if os.IsNotExist(err) {
		return false, &PathNotFound{path}
	}
	return false, err
}
