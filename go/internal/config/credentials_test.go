package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCredentialStore(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	paths := Paths{
		ConfigFile: filepath.Join(tmpDir, "config.toml"),
	}

	tests := []struct {
		name    string
		content string
		want    CredentialStore
	}{
		{
			name: "missing config defaults to unknown",
			want: CredentialStoreUnknown,
		},
		{
			name:    "parses file store",
			content: "cli_auth_credentials_store = \"file\"\n",
			want:    CredentialStoreFile,
		},
		{
			name:    "parses keyring store with comments",
			content: "cli_auth_credentials_store = \"keyring\" # preferred\n",
			want:    CredentialStoreKeyring,
		},
		{
			name:    "parses literal string",
			content: "cli_auth_credentials_store = 'auto'\n",
			want:    CredentialStoreAuto,
		},
		{
			name:    "parses ephemeral store",
			content: "cli_auth_credentials_store = \"ephemeral\"\n",
			want:    CredentialStoreEphemeral,
		},
		{
			name:    "ignores section scoped keys",
			content: "[profiles.work]\ncli_auth_credentials_store = \"file\"\n",
			want:    CredentialStoreUnknown,
		},
		{
			name:    "ignores unsupported values",
			content: "cli_auth_credentials_store = \"mystery\"\n",
			want:    CredentialStoreUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.content == "" {
				_ = os.Remove(paths.ConfigFile)
			} else if err := os.WriteFile(paths.ConfigFile, []byte(tc.content), 0600); err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			got, err := ResolveCredentialStore(paths)
			if err != nil {
				t.Fatalf("ResolveCredentialStore failed: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
