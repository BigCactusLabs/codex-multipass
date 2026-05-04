package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Codex profile compatibility",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()
		report, err := profile.Diagnose(paths)
		if err != nil {
			fail("%v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			if err := writeJSON(os.Stdout, report); err != nil {
				fail("Failed to encode doctor report: %v", err)
			}
			if !report.OK {
				exitFunc(1)
				panic(exitSignal{Code: 1})
			}
			return
		}

		fmt.Println("  Codex Multipass Doctor")
		fmt.Println("  ----------------------------")
		fmt.Printf("  CODEX_DIR        = %s\n", report.CodexHome)
		fmt.Printf("  CREDENTIAL_STORE = %s\n", report.CredentialStore)
		fmt.Printf("  AUTH             = %s\n", formatFileStatus(report.AuthFile))
		fmt.Printf("  PROFILES_DIR     = %s\n", formatFileStatus(report.ProfilesDir))
		if report.ActiveProfile != "" {
			fmt.Printf("  ACTIVE_PROFILE   = %s\n", report.ActiveProfile)
		}

		if len(report.Issues) > 0 {
			fmt.Println("")
			fmt.Println("  Issues")
			for _, issue := range report.Issues {
				fmt.Printf("  - %s: %s\n", issue.Code, issue.Message)
			}
			exitFunc(1)
			panic(exitSignal{Code: 1})
		}

		fmt.Println("")
		fmt.Println("  OK")
	},
}

func formatFileStatus(status profile.DoctorFileStatus) string {
	if !status.Exists {
		return status.Path + " (missing)"
	}
	if status.IsDir {
		return status.Path + " (directory, " + status.Mode + ")"
	}
	return status.Path + " (file, " + status.Mode + ")"
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
