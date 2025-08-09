package sandbox

import (
	"log"
	"os"
	"path/filepath"

	"github.com/OminduD/arch-sandbox/filesystem"
	"github.com/OminduD/arch-sandbox/isolation"
	"github.com/OminduD/arch-sandbox/utils"
)

const (
	sandboxDir  = "/home/user/.arch-sandbox"
	tarballURL  = "https://mirrors.kernel.org/archlinux/iso/latest/archlinux-bootstrap-x86_64.tar.gz"
	tarballPath = sandboxDir + "/archlinux-bootstrap.tar.gz"
)

type Sandbox struct {
	Name       string
	Persist    bool
	BaseDir    string
	RootDir    string
	UpperDir   string
	WorkDir    string
	OverlayDir string
}

func NewSandbox(name string, persist bool) (*Sandbox, error) {
	baseDir := filepath.Join(sandboxDir, name)
	return &Sandbox{
		Name:       name,
		Persist:    persist,
		BaseDir:    baseDir,
		RootDir:    filepath.Join(baseDir, "root"),
		UpperDir:   filepath.Join(baseDir, "upper"),
		WorkDir:    filepath.Join(baseDir, "work"),
		OverlayDir: filepath.Join(baseDir, "overlay"),
	}, nil
}
func (s *Sandbox) Setup() error {
	dirs := []string{s.BaseDir, s.RootDir, s.UpperDir, s.WorkDir, s.OverlayDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	log.Printf("Created directories for sandbox %s", s.Name)
	if err := utils.CheckDependencies(); err != nil {
		return err
	}
	if err := utils.DownloadTarball(tarballURL, tarballPath); err != nil {
		return err
	}
	if err := utils.ExtractTarball(tarballPath, s.RootDir); err != nil {
		return err
	}
	if err := filesystem.SetupOverlay(s.RootDir, s.UpperDir, s.WorkDir, s.OverlayDir); err != nil {
		return err
	}
	return nil
}
func (s *Sandbox) Launch() error {
	log.Printf("Launching sandbox %s", s.Name)
	return isolation.LaunchNspawn(s.OverlayDir, s.Name)
}
func (s *Sandbox) Cleanup() error {
	if s.Persist {
		log.Printf("Persisting sandbox %s", s.Name)
		return nil
	}
	log.Printf("Cleaning up sandbox %s", s.Name)
	if err := filesystem.UnmountOverlay(s.OverlayDir); err != nil {
		return err
	}
	return os.RemoveAll(s.BaseDir)
}
