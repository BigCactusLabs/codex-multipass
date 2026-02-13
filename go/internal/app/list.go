package app

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/spf13/cobra"
)

type ProfileStatus struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Active      bool   `json:"active"`
}

func getFingerprint(path string) (string, error) {
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved profiles",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()
		
		// Get Active Fingerprint
		activeFp, _ := getFingerprint(paths.AuthFile)

		entries, err := os.ReadDir(paths.ProfilesDir)
		if err != nil && !os.IsNotExist(err) {
			fail("Failed to list profiles: %v", err)
		}

		var profiles []ProfileStatus

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".json")
			fullPath := filepath.Join(paths.ProfilesDir, entry.Name())
			fp, _ := getFingerprint(fullPath)
			
			profiles = append(profiles, ProfileStatus{
				Name:        name,
				Fingerprint: fp,
				Active:      (activeFp != "" && fp == activeFp),
			})
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			out := map[string]any{
				"ok":       true,
				"profiles": profiles,
			}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Println("")
			fmt.Println("  Profiles")
			fmt.Println("  ----------------------------")
			for _, p := range profiles {
				short := ""
				if len(p.Fingerprint) >= 12 {
					short = p.Fingerprint[:12]
				}

				if p.Active {
					fmt.Printf("  ▸ %s  %s  active\n", p.Name, short)
				} else {
					fmt.Printf("    %s  %s\n", p.Name, short)
				}
			}
			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
