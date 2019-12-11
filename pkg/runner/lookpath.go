package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Based on LookPath from golang/src/os/exec/lp_unix.go
// If the file starts with a slash, it is an absolute path.
// Otherwise, if it contains a slash, this is unsupported by our version of LookPath.
// Base names are searched inside PATH directories
// The result is always an absolute path. There is no current directory in this context.
func LookPath(file string, path string, root string) (string, error) {
	if file == "" {
		return "", fmt.Errorf("LookPath: filename is empty")
	} else if file[0] == '/' {
		// Absolute path provided
		return file, nil
	} else if strings.Contains(file, "/") {
		return "", fmt.Errorf("LookPath: relative file paths are not supported: \"%s\"", file)
	}

	for _, dir := range filepath.SplitList(path) {
		if dir == "" || dir[0] != '/' {
			// Relative paths in PATH are also ignored
			continue
		}

		path := filepath.Join(root, dir, file)
		if findExecutable(path) == nil {
			return filepath.Join(dir, file), nil
		}
	}
	return "", fmt.Errorf("Executable file \"%s\" not found in path", file)
}

// Copied directly from lp_unix.go. An odd choice of the word 'find'.
func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}
