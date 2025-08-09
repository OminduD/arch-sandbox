package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func CheckDependencies() error {
	for _, cmd := range []string{"systemd-nspawn", "mount", "pacman"} {
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
		cmd := exec.Command("gunzip", "-t", dest)
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
	cmd := exec.Command("gunzip", "-t", dest)
	if err := cmd.Run(); err != nil {
		os.Remove(dest)
		return fmt.Errorf("invalid tarball: %v", err)
	}
	return nil
}

func ExtractTarball(tarballPath, dest string) error {
	log.Println("Extracting tarball")
	file, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
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
