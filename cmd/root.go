package cmd

import (
	"log"
	"os"

	"github.com/OminduD/arch-sandbox/sandbox"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arch-sandbox",
	Short: "Create isolated Arch Linux sandboxes",
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

func init() {
	newCmd.Flags().BoolP("persist", "p", false, "Persist sandbox after exit")
	rootCmd.AddCommand(newCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
