package walker

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Config defines the filtering rules for the walker.
type Config struct {
	Extensions []string
	SkipDirs   []string
}

// DefaultConfig returns the standard configuration for VB6 projects.
func DefaultConfig() Config {
	return Config{
		Extensions: []string{".vbp", ".frm", ".bas", ".cls", ".ctl", ".txt", ".ini", ".cfg", ".md", ".json", ".xml"},
		SkipDirs:   []string{".git", ".svn", "bin", "obj", ".idea", ".vscode", "node_modules"},
	}
}

// Walker handles directory traversal.
type Walker struct {
	config Config
}

// New creates a new Walker with the given config.
func New(config Config) *Walker {
	return &Walker{config: config}
}

// Walk traverses the root directory and returns a list of valid text files.
func (w *Walker) Walk(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if w.shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if w.isValidExtension(path) && w.isProbablyTextFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func (w *Walker) shouldSkipDir(name string) bool {
	// Case-insensitive check for skipped directories
	lowerName := strings.ToLower(name)
	for _, skip := range w.config.SkipDirs {
		if strings.EqualFold(lowerName, skip) {
			return true
		}
	}
	return false
}

func (w *Walker) isValidExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, valid := range w.config.Extensions {
		if ext == valid {
			return true
		}
	}
	// Special case: include files simply named "README" or "LICENSE" without extension?
	// For now, stick to extensions as requested.
	return false
}

// isProbablyTextFile reads the first 1024 bytes to check for NUL bytes.
func (w *Walker) isProbablyTextFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false // Can't read, treat as unsafe
	}
	defer f.Close()

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	// Check for NUL byte, which usually indicates binary
	return !bytes.Contains(buf[:n], []byte{0})
}
