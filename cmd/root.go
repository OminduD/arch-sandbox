package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/OminduD/arch-sandbox/sandbox"
	"github.com/OminduD/arch-sandbox/snapshot"
	"github.com/spf13/cobra"
)

var (
	baseDir string
)

func getDefaultBaseDir() string {
	usr, err := user.Current()
	if err != nil {
		// fallback to environment variable
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, ".arch-sandbox")
		}
		return "/tmp/.arch-sandbox"
	}
	return filepath.Join(usr.HomeDir, ".arch-sandbox")
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", getDefaultBaseDir(), "Base directory for sandboxes")
	rootCmd.AddCommand(snapshotCmd)

	// Add the new command
	newCmd.Flags().BoolP("persist", "p", false, "Persist sandbox after exit")
	newCmd.Flags().StringP("network", "n", "host", "Network mode: host, private, none")
	newCmd.Flags().StringSlice("dns", nil, "Custom DNS servers")
	newCmd.Flags().StringSlice("port", nil, "Port mappings (e.g., host:container)")
	rootCmd.AddCommand(newCmd)

	// Add the install command
	rootCmd.AddCommand(installCmd)
}

var rootCmd = &cobra.Command{
	Use:   "arch-sandbox",
	Short: "Create isolated Arch Linux sandboxes",
}

// Commands for new sandbox
var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new sandbox",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		persist, _ := cmd.Flags().GetBool("persist")
		sb, err := sandbox.NewSandboxWithBaseDir(args[0], persist, baseDir)
		if err != nil { //Error Handling
			log.Fatalf("Failed to create sandbox: %v", err)
		}
		// Create a default SandboxConfig or customize as needed
		config := sandbox.SandboxConfig{}
		if err := sb.Setup(config); err != nil {
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

// Command for snapshot management
var snapshotCmd = &cobra.Command{
	Use:   "snapshot <name> <action> [snapshot-id]",
	Short: "Manage sandbox snapshots",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sandboxPath := filepath.Join(baseDir, args[0])
		action := args[1]
		switch action {
		case "save":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for save action")
			}
			if err := snapshot.SaveSnapshot(sandboxPath, args[2]); err != nil {
				log.Fatalf("Failed to save snapshot: %v", err)
			}
		case "restore":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for restore action")
			}
			if err := snapshot.RestoreSnapshot(sandboxPath, args[2]); err != nil {
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
		sb, err := sandbox.NewSandboxWithBaseDir(args[0], true, baseDir)
		if err != nil {
			log.Fatalf("Failed to load sandbox: %v", err)
		}
		cmdExec := exec.Command("arch-chroot", sb.OverlayDir, "yay", "-S", "--noconfirm", args[1])
		if err := cmdExec.Run(); err != nil {
			log.Fatalf("Failed to install package: %v", err)
		}
	}}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
