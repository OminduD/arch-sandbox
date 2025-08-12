package sandbox

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/OminduD/arch-sandbox/filesystem"
	"github.com/OminduD/arch-sandbox/isolation"
	"github.com/OminduD/arch-sandbox/utils"
	yaml "gopkg.in/yaml.v3"
)

const (
	tarballURL = "https://archive.archlinux.org/iso/2025.07.01/archlinux-bootstrap-2025.07.01-x86_64.tar.zst"
)

type Sandbox struct {
	Name       string
	Persist    bool
	BaseDir    string
	RootDir    string
	UpperDir   string
	WorkDir    string
	OverlayDir string
	TarballURL string
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

func NewSandboxWithBaseDir(name string, persist bool, baseDir string) (*Sandbox, error) {
	sandboxBase := filepath.Join(baseDir, name)
	return &Sandbox{
		Name:       name,
		Persist:    persist,
		BaseDir:    sandboxBase,
		RootDir:    filepath.Join(sandboxBase, "root"),
		UpperDir:   filepath.Join(sandboxBase, "upper"),
		WorkDir:    filepath.Join(sandboxBase, "work"),
		OverlayDir: filepath.Join(sandboxBase, "overlay"),
		TarballURL: tarballURL,
	}, nil
}

type SandboxConfig struct {
	Name     string   `yaml:"name"`
	Persist  bool     `yaml:"persist"`
	Tarball  string   `yaml:"tarball"`
	Packages []string `yaml:"packages"`
	Mounts   []struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	} `yaml:"mounts"`
	Network string `yaml:"network"`
}

func NewSandboxFromConfig(configPath string) (*Sandbox, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg SandboxConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	sb, err := NewSandbox(cfg.Name, cfg.Persist)
	if err != nil {
		return nil, err
	}
	sb.TarballURL = cfg.Tarball
	// Apply mounts and packages in Setup()
	return sb, nil
}
func (s *Sandbox) Setup(cfg SandboxConfig) error {
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

	for _, pkg := range cfg.Packages {
		cmd := exec.Command("arch-chroot", s.OverlayDir, "pacman", "-S", "--noconfirm", pkg)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	for _, mount := range cfg.Mounts {
		cmd := exec.Command("mount", "--bind", mount.Source, filepath.Join(s.OverlayDir, mount.Target))
		if err := cmd.Run(); err != nil {
			return err
		}
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

func (s *Sandbox) InstallAURHelper(helper string) error {
	log.Printf("Installing AUR helper %s", helper)
	cmd := exec.Command("arch-chroot", s.OverlayDir, "/bin/bash", "-c",
		"pacman -S --noconfirm git base-devel && "+
			"useradd -m aur && su aur -c 'git clone https://aur.archlinux.org/"+helper+".git /tmp/"+helper+" && "+
			"cd /tmp/"+helper+" && makepkg -si --noconfirm'")
	return cmd.Run()
}
