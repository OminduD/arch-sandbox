package cmd

//Import packages
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

// rootCmd represents the base command when called without any subcommands
// It serves as the entry point for the CLI application.
var rootCmd = &cobra.Command{
	Use:   "arch-sandbox",
	Short: "Create and manage isolated Arch Linux sandboxes",
	Long: `arch-sandbox is a CLI tool to spin up isolated Arch Linux environments
using overlay filesystems and systemd-nspawn. Ideal for developers,
testers, and enthusiasts who need a clean, disposable environment.`,
}

// newCmd represents the new command
// It allows users to create a new sandbox with various configuration options.
var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new sandbox",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		persist, _ := cmd.Flags().GetBool("persist")
		networkMode, _ := cmd.Flags().GetString("network")
		dns, _ := cmd.Flags().GetStringSlice("dns")
		ports, _ := cmd.Flags().GetStringSlice("port")
		cpuShares, _ := cmd.Flags().GetString("cpu-shares")
		memoryLimit, _ := cmd.Flags().GetString("memory-limit")

		sb, err := sandbox.NewSandboxWithBaseDir(name, persist, baseDir)
		if err != nil {
			log.Fatalf("Failed to create sandbox: %v", err)
		}

		// The config is for features like pre-installing packages from a file, not yet implemented via CLI flags.
		config := sandbox.SandboxConfig{}
		if err := sb.Setup(config); err != nil {
			log.Fatalf("Sandbox setup failed: %v", err)
		}

		if err := sb.Launch(networkMode, dns, ports, cpuShares, memoryLimit); err != nil {
			log.Fatalf("Sandbox launch failed: %v", err)
		}

		// Cleanup is handled after the sandbox session ends.
		if err := sb.Cleanup(); err != nil {
			log.Fatalf("Sandbox cleanup failed: %v", err)
		}
	},
}

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot <name> <action> [snapshot-id]",
	Short: "Manage sandbox snapshots (save, restore)",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sandboxName := args[0]
		action := args[1]
		sandboxPath := filepath.Join(baseDir, sandboxName)

		switch action {
		case "save":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for save action")
			}
			snapshotID := args[2]
			if err := snapshot.SaveSnapshot(sandboxPath, snapshotID); err != nil {
				log.Fatalf("Failed to save snapshot: %v", err)
			}
			log.Printf("Snapshot '%s' saved for sandbox '%s'.\n", snapshotID, sandboxName)
		case "restore":
			if len(args) < 3 {
				log.Fatalf("Missing snapshot-id for restore action")
			}
			snapshotID := args[2]
			if err := snapshot.RestoreSnapshot(sandboxPath, snapshotID); err != nil {
				log.Fatalf("Failed to restore snapshot: %v", err)
			}
			log.Printf("Snapshot '%s' restored for sandbox '%s'.\n", snapshotID, sandboxName)
		default:
			log.Fatalf("Unknown action: %s. Use 'save' or 'restore'.", action)
		}
	},
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <name> <package>",
	Short: "Install a package in a persistent sandbox",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sandboxName := args[0]
		packageName := args[1]
		// Installing implies the sandbox must be persistent.
		sb, err := sandbox.NewSandboxWithBaseDir(sandboxName, true, baseDir)
		if err != nil {
			log.Fatalf("Failed to load sandbox: %v", err)
		}

		// Ensure AUR helper is installed
		if err := sb.InstallAURHelper("yay"); err != nil {
			log.Printf("Could not install AUR helper, proceeding with pacman: %v", err)
		}

		log.Printf("Installing package '%s' in sandbox '%s'...", packageName, sandboxName)
		// Use arch-chroot to run commands inside the sandbox's filesystem
		cmdExec := exec.Command("arch-chroot", sb.OverlayDir, "yay", "-S", "--noconfirm", packageName)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		if err := cmdExec.Run(); err != nil {
			log.Fatalf("Failed to install package: %v", err)
		}
		log.Printf("Successfully installed '%s' in '%s'.", packageName, sandboxName)
	},
}

func init() {
	// Root command persistent flags
	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", getDefaultBaseDir(), "Base directory for sandboxes")

	// `new` command flags
	newCmd.Flags().BoolP("persist", "p", false, "Persist sandbox after exit")
	newCmd.Flags().String("network", "host", "Network mode: host, private, none")
	newCmd.Flags().StringSlice("dns", []string{}, "Custom DNS servers for private network mode")
	newCmd.Flags().StringSlice("port", []string{}, "Port mappings (e.g., host:container)")
	newCmd.Flags().String("cpu-shares", "", "CPU shares (weight)")
	newCmd.Flags().String("memory-limit", "", "Memory limit (e.g., 1G)")

	// Add subcommands to root
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(installCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getDefaultBaseDir() string {
	usr, err := user.Current()
	if err != nil {
		// Fallback to HOME environment variable
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, ".arch-sandbox")
		}
		// Last resort
		return "/tmp/.arch-sandbox"
	}
	return filepath.Join(usr.HomeDir, ".arch-sandbox")
}
