package isolation

import (
	"log"
	"os/exec"
)

func LaunchNspawn(overlayDir, name string) error {
	log.Printf("Launching systemd-nspawn for %s", name)
	cmd := exec.Command("systemd-nspawn", "--directory", overlayDir, "--machine", name, "/bin/bash")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
