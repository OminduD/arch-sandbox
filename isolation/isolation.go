package isolation

import (
	"log"
	"os/exec"
)

func LaunchNspawn(overlayDir, name, networkMode string, dns []string, ports []string, cpuShares, memoryLimit string) error {
	log.Printf("Launching systemd-nspawn for %s", name)
	cmd := exec.Command("systemd-nspawn", "--directory", overlayDir, "--machine", name, "/bin/bash")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	args := []string{"--directory", overlayDir, "--machine", name}
    switch networkMode {
    case "private":
        args = append(args, "--private-network", "--network-veth")
    case "none":
        args = append(args, "--private-network")
    }
    for _, d := range dns {
        args = append(args, "--resolv-conf=off", "--dns="+d)
    }
    for _, p := range ports {
        args = append(args, "--port="+p)
    }
    args = []string{
        "--directory", overlayDir,
        "--machine", name,
        "--cpu-weight", cpuShares,
        "--memory-limit", memoryLimit,
        "/bin/bash",
    }
    cmd = exec.Command("systemd-nspawn", args...)
    cmd.Stdout = log.Writer()
    cmd.Stderr = log.Writer()
    return cmd.Run()
    cmd.Stdout = log.Writer()
    cmd.Stderr = log.Writer()
    return cmd.Run()
	
}
