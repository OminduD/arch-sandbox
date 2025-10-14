package isolation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// LaunchNspawn constructs and executes the systemd-nspawn command to start the container.
func LaunchNspawn(overlayDir, name, networkMode string, dns []string, ports []string, cpuShares, memoryLimit string) error {
	log.Printf("Launching systemd-nspawn for %s", name)

	args := []string{
		"--directory", overlayDir,
		"--machine", name,
	}

	// Configure networking
	switch networkMode {
	case "private":
		args = append(args, "--network-veth")
	case "none":
		args = append(args, "--private-network")
	case "host":
		// This is the default behavior, no extra flag needed.
	default:
		log.Printf("Warning: unknown network mode %q, using systemd-nspawn default (host)", networkMode)
	}

	// Configure DNS if provided
	if len(dns) > 0 {
		args = append(args, "--resolv-conf=off")
		for _, d := range dns {
			args = append(args, "--dns="+d)
		}
	}

	// Configure port mappings
	for _, p := range ports {
		args = append(args, "--port="+p)
	}

	// Configure resource limits
	if cpuShares != "" {
		// systemd-nspawn uses --cpu-weight for shares
		args = append(args, fmt.Sprintf("--cpu-weight=%s", cpuShares))
	}

	if memoryLimit != "" {
		// systemd-nspawn uses --memory-max for limit
		args = append(args, fmt.Sprintf("--memory-max=%s", memoryLimit))
	}

	// The command to run inside the container
	args = append(args, "/bin/bash")

	cmd := exec.Command("systemd-nspawn", args...)

	// Make the session interactive
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Executing: %s", cmd.String())
	return cmd.Run()
}
