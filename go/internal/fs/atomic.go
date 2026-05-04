package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// AtomicWriteJSON writes the given data to a file atomically.
// It creates a temporary file in the same directory, writes to it,
// sets permissions, and then renames it to the target file.
func AtomicWriteJSON(path string, data any, perm os.FileMode) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	tmpClosed := false
	defer func() {
		_ = os.Remove(tmpName)
		if !tmpClosed {
			_ = tmpFile.Close()
		}
	}()

	enc := json.NewEncoder(tmpFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Close before rename to ensure flush
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpClosed = true

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// AtomicCopy copies a file atomically to a target path with specified permissions.
func AtomicCopy(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	dir := filepath.Dir(dst)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	tmpClosed := false
	defer func() {
		_ = os.Remove(tmpName)
		if !tmpClosed {
			_ = tmpFile.Close()
		}
	}()

	if _, err := io.Copy(tmpFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpClosed = true

	if err := os.Rename(tmpName, dst); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LockFile acquires a shared lock on a file path.
// It returns a release function that must be called to unlock.
func Lock(path string) (func(), error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create lock directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	}, nil
}
