package config

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type CredentialStore string

const (
	CredentialStoreUnknown   CredentialStore = ""
	CredentialStoreAuto      CredentialStore = "auto"
	CredentialStoreFile      CredentialStore = "file"
	CredentialStoreKeyring   CredentialStore = "keyring"
	CredentialStoreEphemeral CredentialStore = "ephemeral"
)

var credentialStorePattern = regexp.MustCompile(`^\s*cli_auth_credentials_store\s*=\s*(?:"([^"]+)"|'([^']+)')`)

func ResolveCredentialStore(paths Paths) (CredentialStore, error) {
	file, err := os.Open(paths.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return CredentialStoreUnknown, nil
		}
		return CredentialStoreUnknown, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	inSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inSection = true
			continue
		}
		if inSection {
			continue
		}

		match := credentialStorePattern.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}

		value := strings.ToLower(match[1])
		if value == "" {
			value = strings.ToLower(match[2])
		}

		switch CredentialStore(value) {
		case CredentialStoreAuto, CredentialStoreFile, CredentialStoreKeyring, CredentialStoreEphemeral:
			return CredentialStore(value), nil
		default:
			return CredentialStoreUnknown, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return CredentialStoreUnknown, err
	}

	return CredentialStoreUnknown, nil
}
