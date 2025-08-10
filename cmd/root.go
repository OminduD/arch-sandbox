package cmd

import (
	"log"
	"os"
	"os/exec"

	"path/filepath"

	"github.com/OminduD/arch-sandbox/sandbox"
	"github.com/OminduD/arch-sandbox/snapshot" // Import the snapshot package
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arch-sandbox",
	Short: "Create isolated Arch Linux sandboxes",
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new sandbox",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		persist, _ := cmd.Flags().GetBool("persist")
		sb, err := sandbox.NewSandbox(args[0], persist)
		if err != nil {
			log.Fatalf("Failed to create sandbox: %v", err)
		}
		if err := sb.Setup(); err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
		if err := sb.Launch(); err != nil {
			log.Fatalf("Launch failed: %v", err)
		}
		if err := sb.Cleanup(); err != nil {
			log.Fatalf("Cleanup failed: %v", err)
		}
	},
}
var snapshotCmd = &cobra.Command{
	Use:   "snapshot <name> <action> [snapshot-id]",
	Short: "Manage sandbox snapshots",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sandboxDir := "/path/to/sandbox" // Define sandboxDir or replace with the correct path
		action := args[1]
		switch action {
		case "save":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for save action")
			}
			if err := snapshot.SaveSnapshot(filepath.Join(sandboxDir, args[0]), args[2]); err != nil {
				log.Fatalf("Failed to save snapshot: %v", err)
			}
		case "restore":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for restore action")
			}
			if err := snapshot.RestoreSnapshot(filepath.Join(sandboxDir, args[0]), args[2]); err != nil {
				log.Fatalf("Failed to restore snapshot: %v", err)
			}
		default:
			log.Fatalf("Unknown action: %s", action)
		}
	},
}
var installCmd = &cobra.Command{
	Use:   "install <name> <package>",
	Short: "Install a package in the sandbox",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sb, err := sandbox.NewSandbox(args[0], true)
		if err != nil {
			log.Fatalf("Failed to load sandbox: %v", err)
			cmd := exec.Command("arch-chroot", sb.OverlayDir, "yay", "-S", "--noconfirm", args[1])
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to install package: %v", err)
			}
		}

	}}

func init() {
	newCmd.Flags().BoolP("persist", "p", false, "Persist sandbox after exit")
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("network", "n", "host", "Network mode: host, private, none")
	newCmd.Flags().StringSlice("dns", nil, "Custom DNS servers")
	newCmd.Flags().StringSlice("port", nil, "Port mappings (e.g., host:container)")
	newCmd.Flags().StringP("network", "n", "host", "Network mode: host, private, none")
	newCmd.Flags().StringSlice("dns", nil, "Custom DNS servers")
	newCmd.Flags().StringSlice("port", nil, "Port mappings (e.g., host:container)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
