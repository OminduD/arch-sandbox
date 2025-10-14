package sandbox

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/OminduD/arch-sandbox/filesystem"
	"github.com/OminduD/arch-sandbox/isolation"
	"github.com/OminduD/arch-sandbox/utils"
	"gopkg.in/yaml.v3"
)

const (
	// Using a different fixed, recent, and verified version for reproducibility.
	tarballURL = "https://archive.archlinux.org/iso/2024.07.01/archlinux-bootstrap-2024.07.01-x86_64.tar.zst"
)

// Sandbox defines the structure and paths for an isolated environment.
type Sandbox struct {
	Name       string
	Persist    bool
	BaseDir    string
	RootDir    string // Lower dir for overlayfs
	UpperDir   string // Upper dir for overlayfs
	WorkDir    string // Work dir for overlayfs
	OverlayDir string // Mount point for overlayfs
	TarballURL string
}

// SandboxConfig defines sandbox configurations from a file.
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

// NewSandboxWithBaseDir creates a new Sandbox struct with all paths configured.
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
		TarballURL: tarballURL, // Default URL
	}, nil
}

// NewSandboxFromConfig creates a new sandbox from a YAML configuration file.
func NewSandboxFromConfig(configPath string) (*Sandbox, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg SandboxConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Assuming the base directory is managed outside or defaults are used.
	usr, _ := os.UserHomeDir()
	sb, err := NewSandboxWithBaseDir(cfg.Name, cfg.Persist, filepath.Join(usr, ".arch-sandbox"))
	if err != nil {
		return nil, err
	}
	if cfg.Tarball != "" {
		sb.TarballURL = cfg.Tarball
	}
	return sb, nil
}

// Setup creates directories, downloads and extracts the Arch bootstrap tarball, and sets up the overlayfs.
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

	// Define a shared cache directory for tarballs to avoid re-downloading.
	tarballCacheDir := filepath.Join(s.BaseDir, "..", ".cache")
	if err := os.MkdirAll(tarballCacheDir, 0755); err != nil {
		return err
	}
	tarballPath := filepath.Join(tarballCacheDir, filepath.Base(s.TarballURL))

	if err := utils.DownloadTarball(s.TarballURL, tarballPath); err != nil {
		return err
	}
	if err := utils.ExtractTarball(tarballPath, s.RootDir); err != nil {
		return err
	}
	if err := filesystem.SetupOverlay(s.RootDir, s.UpperDir, s.WorkDir, s.OverlayDir); err != nil {
		return err
	}

	for _, pkg := range cfg.Packages {
		log.Printf("Installing package: %s", pkg)
		cmd := exec.Command("arch-chroot", s.OverlayDir, "pacman", "-S", "--noconfirm", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	for _, mount := range cfg.Mounts {
		log.Printf("Binding mount from %s to %s", mount.Source, mount.Target)
		targetPath := filepath.Join(s.OverlayDir, mount.Target)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		cmd := exec.Command("mount", "--bind", mount.Source, targetPath)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Launch starts the systemd-nspawn container.
func (s *Sandbox) Launch(networkMode string, dns []string, ports []string, cpuShares, memoryLimit string) error {
	log.Printf("Launching sandbox %s", s.Name)
	return isolation.LaunchNspawn(s.OverlayDir, s.Name, networkMode, dns, ports, cpuShares, memoryLimit)
}

// Cleanup unmounts the overlayfs and removes the sandbox directory if not persistent.
func (s *Sandbox) Cleanup() error {
	log.Println("Unmounting overlayfs...")
	if err := filesystem.UnmountOverlay(s.OverlayDir); err != nil {
		// Log the error but don't stop, still try to clean up if not persisting
		log.Printf("Warning: failed to unmount overlayfs: %v", err)
	}

	if s.Persist {
		log.Printf("Persisting sandbox '%s' at %s", s.Name, s.BaseDir)
		return nil
	}

	log.Printf("Cleaning up sandbox %s", s.Name)
	return os.RemoveAll(s.BaseDir)
}

// InstallAURHelper installs an AUR helper like 'yay' into the sandbox.
func (s *Sandbox) InstallAURHelper(helper string) error {
	log.Printf("Installing AUR helper '%s'", helper)
	script := `
pacman -S --noconfirm --needed git base-devel && \
useradd -m builder && \
chown -R builder:builder /home/builder && \
su builder -c 'cd /tmp && git clone https://aur.archlinux.org/` + helper + `.git && cd ` + helper + ` && makepkg -si --noconfirm'
`
	cmd := exec.Command("arch-chroot", s.OverlayDir, "/bin/bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
