package config

import (
	"os"
	"path/filepath"
)

// Paths holds the resolved paths for the application
type Paths struct {
	CodexHome   string `json:"codex_home,omitempty"` // The env var value, if set
	CodexDir    string `json:"codex_dir"`
	AuthFile    string `json:"auth"`
	ProfilesDir string `json:"profiles_dir"`
	ActiveFile  string `json:"-"`
}

// ResolvePaths determines the runtime paths based on environment variables and defaults.
func ResolvePaths() Paths {
	homeDir, _ := os.UserHomeDir()

	// Default: ~/.codex
	codexDir := filepath.Join(homeDir, ".codex")

	// Override if CODEX_HOME is set
	envHome := os.Getenv("CODEX_HOME")
	if envHome != "" {
		codexDir = envHome
	}

	return Paths{
		CodexHome:   envHome,
		CodexDir:    codexDir,
		AuthFile:    filepath.Join(codexDir, "auth.json"),
		ProfilesDir: filepath.Join(codexDir, "profiles"),
		ActiveFile:  filepath.Join(codexDir, ".codex-mp-active"),
	}
}
