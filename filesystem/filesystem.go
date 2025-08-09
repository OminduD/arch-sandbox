package filesystem

import (
	"log"
	"os/exec"
)

func SetupOverlay(lowerDir, upperDir, workDir, overlayDir string) error {
	log.Println("Setting up overlayfs")
	cmd := exec.Command("mount", "-t", "overlay", "overlay",
		"-o", "lowerdir="+lowerDir+",upperdir="+upperDir+",workdir="+workDir, overlayDir)
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Overlayfs mounted")
	return nil
}
func UnmountOverlay(overlayDir string) error {
	log.Println("Unmounting overlayfs")
	cmd := exec.Command("umount", overlayDir)
	return cmd.Run()
}
