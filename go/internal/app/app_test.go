package app

import (
	"testing"
)

func runAndCaptureExit(t *testing.T, fn func()) (exitCode int) {
	t.Helper()

	originalExit := exitFunc
	exitFunc = func(code int) {
		panic(exitSignal{Code: code})
	}
	defer func() {
		exitFunc = originalExit
	}()

	exitCode = -1
	defer func() {
		if r := recover(); r != nil {
			sig, ok := r.(exitSignal)
			if !ok {
				panic(r)
			}
			exitCode = sig.Code
		}
	}()

	fn()
	return exitCode
}

func TestFailExitsWithCode1(t *testing.T) {
	code := runAndCaptureExit(t, func() {
		fail("boom")
	})

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestSaveWithoutAuthExits(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	rootCmd.SetArgs([]string{"save", "my-profile"})

	code := runAndCaptureExit(t, func() {
		_ = rootCmd.Execute()
	})

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}
