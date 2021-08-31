package fs

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	dirPerms = 0700
)

// DefaultVegaDir returns the location to Vega config files and data files:
// binary is in /usr/bin/ -> look for /etc/vega/config.toml
// binary is in /usr/local/vega/bin/ -> look for /usr/local/vega/etc/config.toml
// binary is in /usr/local/bin/ -> look for /usr/local/etc/vega/config.toml
// otherwise, look for $HOME/.vega/config.toml
func DefaultVegaDir() string {
	if runtime.GOOS == "windows" {
		// shortcut for windows
		p, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(p, ".vega")
		}
	}

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
			return os.MkdirAll(path, dirPerms)
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
		return false, nil
	}
	return false, err
}

// FileExists similar to PathExists, but ensures the path is to a file, not a
// directory.
func FileExists(path string) (bool, error) {
	fs, err := os.Stat(path)
	if err == nil {
		return !fs.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func WriteFile(content []byte, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}
