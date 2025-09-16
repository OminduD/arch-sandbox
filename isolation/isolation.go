package isolation

import (
	"log"
	"os/exec"
)

func LaunchNspawn(overlayDir, name, networkMode string, dns []string, ports []string, cpuShares, memoryLimit string) error {
	log.Printf("Launching systemd-nspawn for %s", name)

	args := []string{"--directory", overlayDir, "--machine", name}

	// Configure network mode
	switch networkMode {
	case "private":
		args = append(args, "--private-network", "--network-veth")
	case "none":
		args = append(args, "--private-network")
		// default "host" mode needs no additional flags
	}

	// Add DNS configuration if provided
	for _, d := range dns {
		args = append(args, "--resolv-conf=off", "--dns="+d)
	}

	// Add port mappings if provided
	for _, p := range ports {
		args = append(args, "--port="+p)
	}

	// Add resource limits if provided
	if cpuShares != "" {
		args = append(args, "--cpu-weight", cpuShares)
	}
	if memoryLimit != "" {
		args = append(args, "--memory-limit", memoryLimit)
	}

	// Add the command to run
	args = append(args, "/bin/bash")

	cmd := exec.Command("systemd-nspawn", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
