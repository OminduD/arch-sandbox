package snapshot

import (
    "os"
    "os/exec"
    "path/filepath"
)

func SaveSnapshot(sandboxDir, snapshotName string) error {
    snapshotPath := filepath.Join(sandboxDir, "snapshots", snapshotName+".tar.zst")
    if err := os.MkdirAll(filepath.Dir(snapshotPath), 0755); err != nil {
        return err
    }
    cmd := exec.Command("tar", "-C", filepath.Join(sandboxDir, "upper"), "--zstd", "-cf", snapshotPath, ".")
    return cmd.Run()
}

func RestoreSnapshot(sandboxDir, snapshotName string) error {
    upperDir := filepath.Join(sandboxDir, "upper")
    os.RemoveAll(upperDir)
    os.MkdirAll(upperDir, 0755)
    snapshotPath := filepath.Join(sandboxDir, "snapshots", snapshotName+".tar.zst")
    cmd := exec.Command("tar", "-C", upperDir, "--zstd", "-xf", snapshotPath)
    return cmd.Run()
}