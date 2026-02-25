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
		ActiveFile:  filepath.Join(tmpDir, ".codex-mp-active"),
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

func TestUseSameProfileDoesNotSyncIntoProfile(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"stable"}`), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}
	if _, err := Save("same", paths); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Simulate auth being changed outside codex-mp while same profile remains active.
	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"temp"}`), 0600); err != nil {
		t.Fatalf("failed to mutate auth file: %v", err)
	}

	if err := Use("same", paths); err != nil {
		t.Fatalf("use failed: %v", err)
	}

	profileRaw, err := os.ReadFile(filepath.Join(paths.ProfilesDir, "same.json"))
	if err != nil {
		t.Fatalf("failed to read profile file: %v", err)
	}
	if string(profileRaw) != `{"token":"stable"}` {
		t.Fatalf("expected profile to remain stable, got %s", string(profileRaw))
	}

	authRaw, err := os.ReadFile(paths.AuthFile)
	if err != nil {
		t.Fatalf("failed to read auth file: %v", err)
	}
	if string(authRaw) != `{"token":"stable"}` {
		t.Fatalf("expected auth to be restored from profile, got %s", string(authRaw))
	}
}

func TestUseSyncsActiveProfileBeforeSwitch(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"a-v1"}`), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}
	if _, err := Save("a", paths); err != nil {
		t.Fatalf("failed to save profile a: %v", err)
	}

	// Simulate Codex refreshing tokens for the active profile in auth.json.
	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"a-v2"}`), 0600); err != nil {
		t.Fatalf("failed to update auth file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(paths.ProfilesDir, "b.json"), []byte(`{"token":"b-v1"}`), 0600); err != nil {
		t.Fatalf("failed to write profile b: %v", err)
	}

	if err := Use("b", paths); err != nil {
		t.Fatalf("use failed: %v", err)
	}

	aRaw, err := os.ReadFile(filepath.Join(paths.ProfilesDir, "a.json"))
	if err != nil {
		t.Fatalf("failed to read synced profile a: %v", err)
	}
	if string(aRaw) != `{"token":"a-v2"}` {
		t.Fatalf("expected profile a to be synced with refreshed auth, got %s", string(aRaw))
	}

	activeRaw, err := os.ReadFile(paths.ActiveFile)
	if err != nil {
		t.Fatalf("failed to read active marker: %v", err)
	}
	if got := string(activeRaw); got != "b\n" {
		t.Fatalf("expected active marker to be b, got %q", got)
	}
}

func TestListUsesActiveMarkerWhenFingerprintChanges(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"a-v1"}`), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}
	if _, err := Save("a", paths); err != nil {
		t.Fatalf("failed to save profile a: %v", err)
	}

	// Auth rotates and no longer matches the saved profile fingerprint.
	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"a-v2"}`), 0600); err != nil {
		t.Fatalf("failed to rotate auth file: %v", err)
	}

	profiles, err := List(paths)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	if !profiles[0].Active {
		t.Fatalf("expected profile a to still be marked active")
	}
}

func TestRenameAndDeleteUpdateActiveMarker(t *testing.T) {
	paths, cleanup := setupTest(t)
	defer cleanup()

	if err := os.WriteFile(paths.AuthFile, []byte(`{"token":"x"}`), 0600); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}
	if _, err := Save("work", paths); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := Rename("work", "office", paths); err != nil {
		t.Fatalf("rename failed: %v", err)
	}

	activeRaw, err := os.ReadFile(paths.ActiveFile)
	if err != nil {
		t.Fatalf("failed to read active marker after rename: %v", err)
	}
	if got := string(activeRaw); got != "office\n" {
		t.Fatalf("expected active marker office after rename, got %q", got)
	}

	if err := Delete("office", paths); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := os.Stat(paths.ActiveFile); !os.IsNotExist(err) {
		t.Fatalf("expected active marker to be removed after deleting active profile")
	}
}
