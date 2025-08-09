package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func CheckDependencies() error {
	for _, cmd := range []string{"systemd-nspawn", "mount", "pacman", "zstd"} { // add zstd
		if _, err := exec.LookPath(cmd); err != nil {
			return err
		}
	}
	return nil
}

func DownloadTarball(url, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		// Verify existing tarball
		log.Println("Verifying existing tarball")
		cmd := exec.Command("zstd", "-t", dest) // use zstd for .zst
		if err := cmd.Run(); err == nil {
			log.Println("Tarball already exists and is valid")
			return nil
		}
		log.Println("Existing tarball is invalid, redownloading")
		os.Remove(dest)
	}
	log.Println("Downloading Arch Linux tarball")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	// Verify downloaded tarball
	log.Println("Verifying downloaded tarball")
	cmd := exec.Command("zstd", "-t", dest) // use zstd for .zst
	if err := cmd.Run(); err != nil {
		os.Remove(dest)
		return fmt.Errorf("invalid tarball: %v", err)
	}
	return nil
}

// Replace ExtractTarball to handle .tar.zst
func ExtractTarball(tarballPath, dest string) error {
	log.Println("Extracting tarball")

	// Define and initialize cmd for tar extraction
	cmd := exec.Command("tar", "--use-compress-program=zstd", "-xf", tarballPath, "-C", dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tar extraction: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	// Set root ownership and correct permissions
	if err := os.Chown(dest, 0, 0); err != nil { // 0 is root's UID/GID
		return fmt.Errorf("failed to chown root dir: %v", err)
	}
	if err := filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := os.Chown(path, 0, 0); err != nil { // Set root ownership
			return err
		}
		if info.Mode().IsRegular() && (path == filepath.Join(dest, "bin/bash") || info.Mode()&0111 != 0) {
			// Ensure executable files have execute permissions
			return os.Chmod(path, info.Mode()|0111)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to set permissions: %v", err)
	}

	// Handle Arch bootstrap tarball structure (root.x86_64)
	extractedRoot := filepath.Join(dest, "root.x86_64")
	if _, err := os.Stat(extractedRoot); err == nil {
		if err := moveContents(extractedRoot, dest); err != nil {
			return err
		}
		os.RemoveAll(extractedRoot)
	}
	log.Println("Tarball extracted")
	return nil
}
func moveContents(src, dst string) error {
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err := os.Rename(filepath.Join(src, name), filepath.Join(dst, name)); err != nil {
			return err
		}
	}
	return nil
}
