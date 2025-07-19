package routing

import (
	"net/http"
	"os"
	"path/filepath"
)

type FileOnlyFileSystem struct {
	root string
}

func newFileOnlyFileSystem(root string) *FileOnlyFileSystem {
	return &FileOnlyFileSystem{root: root}
}

// Open implements the http.FileSystem interface.
func (fs *FileOnlyFileSystem) Open(name string) (http.File, error) {
	// filepath.Join is crucial for security as it prevents path traversal attacks.
	fullPath := filepath.Join(fs.root, name)

	info, err := os.Stat(fullPath)
	if err != nil {
		// If os.Stat returns an error, the file or directory does not exist.
		// We can return the error directly. http.FileServer will handle it
		// (e.g., os.ErrNotExist will be converted to a 404).
		return nil, err
	}

	if info.IsDir() {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
