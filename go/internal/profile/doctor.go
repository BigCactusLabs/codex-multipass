package profile

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
)

type DoctorFileStatus struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
	IsDir  bool   `json:"is_dir"`
	Mode   string `json:"mode,omitempty"`
}

type DoctorIssue struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type DoctorReport struct {
	OK              bool             `json:"ok"`
	CodexHome       string           `json:"codex_home"`
	CredentialStore string           `json:"credential_store"`
	AuthFile        DoctorFileStatus `json:"auth_file"`
	ProfilesDir     DoctorFileStatus `json:"profiles_dir"`
	ActiveProfile   string           `json:"active_profile,omitempty"`
	Issues          []DoctorIssue    `json:"issues"`
}

func Diagnose(paths config.Paths) (DoctorReport, error) {
	store, err := config.ResolveCredentialStore(paths)
	if err != nil {
		return DoctorReport{}, fmt.Errorf("failed to resolve credential store: %w", err)
	}

	activeProfile, err := readActiveProfile(paths)
	if err != nil {
		return DoctorReport{}, err
	}

	report := DoctorReport{
		OK:              true,
		CodexHome:       paths.CodexDir,
		CredentialStore: credentialStoreLabel(store),
		AuthFile:        statPath(paths.AuthFile),
		ProfilesDir:     statPath(paths.ProfilesDir),
		ActiveProfile:   activeProfile,
		Issues:          []DoctorIssue{},
	}

	switch store {
	case config.CredentialStoreKeyring, config.CredentialStoreAuto, config.CredentialStoreEphemeral:
		report.addIssue("unsupported_credential_store", fmt.Sprintf("Codex is configured for %q credential storage. codex-mp currently manages only file-backed auth.json profiles; set cli_auth_credentials_store = \"file\" in %s and sign in again.", store, paths.ConfigFile))
	}

	if !report.AuthFile.Exists {
		report.addIssue("missing_auth_file", fmt.Sprintf("Missing auth file at %s. Run 'codex login' after setting cli_auth_credentials_store = \"file\".", paths.AuthFile))
	} else if report.AuthFile.IsDir {
		report.addIssue("auth_file_is_directory", fmt.Sprintf("Auth path is a directory, not a file: %s", paths.AuthFile))
	}

	if report.ProfilesDir.Exists && !report.ProfilesDir.IsDir {
		report.addIssue("profiles_path_is_file", fmt.Sprintf("Profiles path exists but is not a directory: %s", paths.ProfilesDir))
	}

	report.OK = len(report.Issues) == 0
	return report, nil
}

func (r *DoctorReport) addIssue(code, message string) {
	r.Issues = append(r.Issues, DoctorIssue{
		Code:    code,
		Message: message,
	})
}

func credentialStoreLabel(store config.CredentialStore) string {
	if store == config.CredentialStoreUnknown {
		return "file"
	}
	return string(store)
}

func statPath(path string) DoctorFileStatus {
	status := DoctorFileStatus{Path: path}
	info, err := os.Stat(path)
	if err != nil {
		return status
	}
	status.Exists = true
	status.IsDir = info.IsDir()
	status.Mode = fmt.Sprintf("%03o", info.Mode().Perm())
	return status
}
