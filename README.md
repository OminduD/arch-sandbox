# Arch Sandbox

**Arch Sandbox** is a CLI tool to create and launch isolated Arch Linux sandboxes using OverlayFS and systemd-nspawn.  
It downloads the official Arch Linux bootstrap tarball (`.tar.zst`), sets up a writable overlay, and launches a containerized shell.

---

## Features

- Automated download and extraction of the latest Arch Linux bootstrap tarball (`.tar.zst`)
- OverlayFS for writable sandboxes
- Launches containers with `systemd-nspawn`
- Optional persistence of sandbox state
- CLI interface with [Cobra](https://github.com/spf13/cobra)

---

## Requirements

- Linux (with OverlayFS support)
- `systemd-nspawn`
- `mount`
- `pacman`
- `zstd`

Install dependencies on Arch Linux:
```sh
sudo pacman -S systemd zstd
```

---

## Build

```sh
git clone https://github.com/OminduD/arch-sandbox.git
cd arch-sandbox
go build -o arch-sandbox
```

---

## Usage

### Create and launch a new sandbox

```sh
sudo ./arch-sandbox new <sandbox-name>
```

Example:
```sh
sudo ./arch-sandbox new testsandbox
```

### Persist sandbox after exit

```sh
sudo ./arch-sandbox new testsandbox --persist
```

---

## How it works

1. **Downloads** the Arch Linux bootstrap tarball (`.tar.zst`)
2. **Extracts** it to a sandbox directory
3. **Sets up OverlayFS** for writable root
4. **Launches** a shell in the sandbox using `systemd-nspawn`
5. **Cleans up** the sandbox unless `--persist` is used

---

## Troubleshooting

- **Permission denied on `/bin/bash`**:  
  Ensure you run the tool as `root` (`sudo`).  
  The extracted files must be owned by root and have correct permissions.

- **OverlayFS errors**:  
  Make sure your kernel supports OverlayFS and the upper/work directories are empty and owned by root.

---

## Project Structure

```
.
├── cmd/         # CLI entrypoint (Cobra)
├── filesystem/  # OverlayFS setup/teardown
├── isolation/   # systemd-nspawn launch logic
├── sandbox/     # Sandbox orchestration
├── snapshot/    # (Placeholder) Snapshot logic
├── utils/       # Download/extract helpers
├── main.go
└── go.mod
```

---

## License

MIT