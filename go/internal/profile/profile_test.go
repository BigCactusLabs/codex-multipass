package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
)

func setupTest(t *testing.T) (config.Paths, func()) {
	tmpDir, err := os.MkdirTemp("", "codex-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	profilesDir := filepath.Join(tmpDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0700); err != nil {
		t.Fatalf("failed to create profiles dir: %v", err)
	}

	paths := config.Paths{
		CodexDir:    tmpDir,
		AuthFile:    filepath.Join(tmpDir, "auth.json"),
		ProfilesDir: profilesDir,
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return paths, cleanup
}

func TestSaveAndList(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	// 1. Create dummy auth file
	authContent := `{"token": "test-token"}`
	if err := os.WriteFile(paths.AuthFile, []byte(authContent), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}

	// 2. Save profile
	name := "test-profile"
	_, err := Save(name, paths)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 3. List profiles
	profiles, err := List(paths)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}

	if profiles[0].Name != name {
		t.Errorf("expected name %s, got %s", name, profiles[0].Name)
	}

	if !profiles[0].Active {
		t.Errorf("expected profile to be active")
	}
}

func TestUseAndRename(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	// 1. Setup two profiles
	auth1 := `{"id": 1}`
	auth2 := `{"id": 2}`

	p1Path := filepath.Join(paths.ProfilesDir, "p1.json")
	p2Path := filepath.Join(paths.ProfilesDir, "p2.json")

	os.WriteFile(p1Path, []byte(auth1), 0600)
	os.WriteFile(p2Path, []byte(auth2), 0600)

	// 2. Use p1
	err := Use("p1", paths)
	if err != nil {
		t.Fatalf("Use failed: %v", err)
	}

	activeFp, _ := GetFingerprint(paths.AuthFile)
	p1Fp, _ := GetFingerprint(p1Path)
	if activeFp != p1Fp {
		t.Errorf("expected p1 to be active")
	}

	// 3. Rename p2 to p3
	err = Rename("p2", "p3", paths)
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(paths.ProfilesDir, "p3.json")); os.IsNotExist(err) {
		t.Errorf("p3.json should exist")
	}
	if _, err := os.Stat(p2Path); !os.IsNotExist(err) {
		t.Errorf("p2.json should not exist")
	}
}

func TestDelete(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	pPath := filepath.Join(paths.ProfilesDir, "gone.json")
	os.WriteFile(pPath, []byte("{}"), 0600)

	err := Delete("gone", paths)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := os.Stat(pPath); !os.IsNotExist(err) {
		t.Errorf("file should be gone")
	}
}
