package profile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/fs"
)

// EnsureInitialized ensures that the profiles directory exists and has correct permissions.
func EnsureInitialized(paths config.Paths) error {
	if err := os.MkdirAll(paths.ProfilesDir, 0700); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	// Enforce 0700 on CodexDir and ProfilesDir explicitly
	if err := os.Chmod(paths.CodexDir, 0700); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", paths.CodexDir, err)
	}
	if err := os.Chmod(paths.ProfilesDir, 0700); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", paths.ProfilesDir, err)
	}

	return nil
}

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// ProfileStatus represents the state of a profile
type ProfileStatus struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Active      bool   `json:"active"`
}

// ValidateName checks if the profile name is valid
func ValidateName(name string) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid profile name: %s (allowed: A-Z a-z 0-9 . _ -)", name)
	}
	return nil
}

// GetFingerprint returns the SHA256 fingerprint of a file
func GetFingerprint(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// withLock executes the given function with a file lock
func withLock(paths config.Paths, action func() error) error {
	if err := EnsureInitialized(paths); err != nil {
		return err
	}

	lockPath := filepath.Join(paths.CodexDir, ".codex-mp.lock")
	unlock, err := fs.Lock(lockPath)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer unlock()

	return action()
}

// Save saves the current auth as a profile
func Save(name string, paths config.Paths) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	profilePath := filepath.Join(paths.ProfilesDir, name+".json")

	err := withLock(paths, func() error {
		// Check Auth Existence INSIDE lock
		if _, err := os.Stat(paths.AuthFile); os.IsNotExist(err) {
			return fmt.Errorf("missing auth file: %s. Hint: run 'codex login' first", paths.AuthFile)
		}

		// Atomic Copy
		if err := fs.AtomicCopy(paths.AuthFile, profilePath, 0600); err != nil {
			return fmt.Errorf("failed to save profile: %w", err)
		}
		return nil
	})

	return profilePath, err
}

// Use switches to a saved profile
func Use(name string, paths config.Paths) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	profilePath := filepath.Join(paths.ProfilesDir, name+".json")

	return withLock(paths, func() error {
		// Check Profile Existence INSIDE lock
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			return fmt.Errorf("profile not found: %s", name)
		}

		// Atomic Copy
		if err := fs.AtomicCopy(profilePath, paths.AuthFile, 0600); err != nil {
			return fmt.Errorf("failed to switch profile: %w", err)
		}
		return nil
	})
}

// Delete removes a profile
func Delete(name string, paths config.Paths) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	profilePath := filepath.Join(paths.ProfilesDir, name+".json")

	return withLock(paths, func() error {
		// Check Existence INSIDE lock
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			return fmt.Errorf("profile not found: %s", name)
		}

		if err := os.Remove(profilePath); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}
		return nil
	})
}

// Rename renames a profile
func Rename(oldName, newName string, paths config.Paths) error {
	if err := ValidateName(oldName); err != nil {
		return err
	}
	if err := ValidateName(newName); err != nil {
		return err
	}

	oldPath := filepath.Join(paths.ProfilesDir, oldName+".json")
	newPath := filepath.Join(paths.ProfilesDir, newName+".json")

	return withLock(paths, func() error {
		// Checks INSIDE lock
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			return fmt.Errorf("profile not found: %s", oldName)
		}
		if _, err := os.Stat(newPath); err == nil {
			return fmt.Errorf("profile already exists: %s", newName)
		}

		// Ensure permissions before rename if possible, or after
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to rename profile: %w", err)
		}

		if err := os.Chmod(newPath, 0600); err != nil {
			return fmt.Errorf("failed to set permissions on renamed profile: %w", err)
		}
		return nil
	})
}

// List returns all profiles
func List(paths config.Paths) ([]ProfileStatus, error) {
	var profiles []ProfileStatus

	err := withLock(paths, func() error {
		// active fingerprint (read inside lock)
		activeFp, _ := GetFingerprint(paths.AuthFile)

		entries, err := os.ReadDir(paths.ProfilesDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil // Return empty list
			}
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".json")
			fullPath := filepath.Join(paths.ProfilesDir, entry.Name())

			fp, err := GetFingerprint(fullPath)
			if err != nil {
				// Profile might have been deleted concurrently, skip it
				continue
			}

			profiles = append(profiles, ProfileStatus{
				Name:        name,
				Fingerprint: fp,
				Active:      (activeFp != "" && fp == activeFp),
			})
		}
		return nil
	})

	return profiles, err
}
