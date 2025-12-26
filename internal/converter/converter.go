package converter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// TransformToUTF8 converts GBK encoded data to UTF-8.
func TransformToUTF8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

// TransformToGBK converts UTF-8 encoded data to GBK.
func TransformToGBK(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}

// ConvertFile converts the file at path.
// If toUTF8 is true, converts GBK -> UTF-8.
// If toUTF8 is false, converts UTF-8 -> GBK.
func ConvertFile(path string, toUTF8 bool) error {
	// Read source
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var newContent []byte
	if toUTF8 {
		// GBK -> UTF-8
		newContent, err = TransformToUTF8(content)
	} else {
		// UTF-8 -> GBK
		newContent, err = TransformToGBK(content)
	}

	if err != nil {
		return fmt.Errorf("transformation failed: %w", err)
	}

	// Safe Write
	return safeWrite(path, newContent)
}

func safeWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	tmpPath := filepath.Join(dir, name+".tmp")
	// bakPath := filepath.Join(dir, name+".bak")

	// 1. Write to tmp
	// Use 0666 for permissions, respecting umask, or copy original perms?
	// Copy original perms is better.
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}
	perm := info.Mode().Perm()

	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return fmt.Errorf("write tmp failed: %w", err)
	}

	// 2. Rename tmp to original (Atomic replace)
	if err := os.Rename(tmpPath, path); err != nil {
		// Clean up tmp on failure
		os.Remove(tmpPath)
		return fmt.Errorf("replace failed: %w", err)
	}

	return nil
}
