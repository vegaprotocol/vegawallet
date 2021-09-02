package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// EnsureDir will make sure a directory exists or is created at a given filesystem path.
func EnsureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModeDir|0700)
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
	fileInfo, err := os.Stat(path)
	if err == nil {
		return !fileInfo.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadFile(path string) ([]byte, error) {
	dir, fileName := filepath.Split(path)
	if len(dir) == 0 {
		dir = "."
	}

	buf, err := fs.ReadFile(os.DirFS(dir), fileName)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file at %s: %w", path, err)
	}

	return buf, nil
}

func WriteFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create file at %s: %w", path, err)
	}
	defer f.Close()

	err = f.Chmod(0600)
	if err != nil {
		return fmt.Errorf("couldn't change file mode at %s: %w", path, err)
	}

	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("couldn't write file at %s: %w", path, err)
	}

	return nil
}
