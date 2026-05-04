package profile

import (
	"fmt"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
)

func ensureFileBackedAuthSupported(paths config.Paths) error {
	store, err := config.ResolveCredentialStore(paths)
	if err != nil {
		return nil
	}
	switch store {
	case config.CredentialStoreUnknown, config.CredentialStoreFile:
		return nil
	}

	return fmt.Errorf("codex is configured for %q credential storage (%s). codex-mp currently manages only file-backed auth.json profiles. Set cli_auth_credentials_store = \"file\" and sign in again", store, paths.ConfigFile)
}

func missingAuthFileError(paths config.Paths) error {
	store, err := config.ResolveCredentialStore(paths)
	if err == nil {
		switch store {
		case config.CredentialStoreKeyring, config.CredentialStoreAuto, config.CredentialStoreEphemeral:
			return ensureFileBackedAuthSupported(paths)
		}
	}

	return fmt.Errorf("missing auth file: %s. Run 'codex login' first. If 'codex login status' says you are already logged in, Codex may be using the OS credential store instead of auth.json. codex-mp only manages file-backed auth.json profiles, so set cli_auth_credentials_store = \"file\" in %s and sign in again", paths.AuthFile, paths.ConfigFile)
}
