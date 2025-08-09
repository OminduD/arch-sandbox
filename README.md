# Arch-Sandbox 🏝️

**Create and manage isolated Arch Linux sandboxes with ease**

`arch-sandbox` is a command-line tool that spins up isolated Arch Linux environments using overlay filesystems and `systemd-nspawn`. Ideal for developers, system administrators, and Linux enthusiasts who want to test, experiment, or develop in a safe, isolated environment.

## ✨ Features
- **Isolated Environments**: Run Arch Linux sandboxes without affecting your host system.
- **Overlay Filesystem**: Keep changes separate using overlayfs.
- **Persistent or Ephemeral**: Choose to keep or discard your sandbox after use.
- **Simple CLI**: Powered by Cobra for a user-friendly command-line experience.

## 📋 Prerequisites
Ensure you have:
- **Go** (1.16 or later) for building.
- **systemd-nspawn** for containerization.
- **mount** for overlay filesystem operations.
- **pacman** for Arch Linux package management.
- **zstd** for `.tar.zst` archives.

Install on Arch Linux:
```bash
sudo pacman -S go systemd zstd
```

## 🚀 Installation
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/OminduD/arch-sandbox.git
   cd arch-sandbox
   ```
2. **Build the Tool**:
   ```bash
   go build -o arch-sandbox
   ```
3. **(Optional) Install Globally**:
   ```bash
   sudo mv arch-sandbox /usr/local/bin/
   ```

## 🛠️ Usage
Create and manage sandboxes with a simple command:

### Create a New Sandbox
```bash
arch-sandbox new <name> [--persist]
```
- `<name>`: Name your sandbox (e.g., `testbox`).
- `--persist` or `-p`: Keep the sandbox after exiting (default: cleanup).

**Examples**:
- Disposable sandbox:
  ```bash
  arch-sandbox new testbox
  ```
- Persistent sandbox:
  ```bash
  arch-sandbox new testbox -p
  ```

### Sandbox Creation Process
The tool follows these steps to create a sandbox:

```
Start
  ↓
[Create Directories] → [Check Dependencies]
                          ↓
                       [Download Tarball]
                          ↓
                       [Extract Tarball]
                          ↓
                       [Mount Overlayfs]
                          ↓
                       [Launch systemd-nspawn]
                          ↓
                       [Cleanup (if not persistent)]
```

**Visualized Process**:
> **Note**: Add a demo GIF to your repo’s assets to see this in action! *(Placeholder: Record running `arch-sandbox new testbox` with `asciinema` or `ffmpeg`)*.

**Estimated Step Durations** (hypothetical):
> **Graph Placeholder**: Add a bar chart image to your repo (e.g., `assets/setup-durations.png`) showing estimated times for each step. Example:
> - Create Directories: ~1s
> - Check Dependencies: ~0.5s
> - Download Tarball: ~30s (depends on network)
> - Extract Tarball: ~10s
> - Mount Overlayfs: ~2s
> - Launch systemd-nspawn: ~5s
>
> *To create this chart, use Chart.js or a graphing tool and host the image in your repo.*

### Inside the Sandbox
You’ll enter a `/bin/bash` shell where you can:
- Install packages with `pacman`.
- Test scripts or configurations.
- Exit with `exit` or `Ctrl+D`.

### Cleanup
- **Disposable**: Automatically deleted on exit.
- **Persistent**: Kept in `~/.arch-sandbox/<name>`. Delete manually:
  ```bash
  rm -rf ~/.arch-sandbox/<name>
  ```

## ⚠️ Notes
- Run as `root` or with `sudo` for `systemd-nspawn` and `mount` operations.
- Internet access is required for tarball download.
- Tarball source: `https://archive.archlinux.org/iso/2025.07.01/archlinux-bootstrap-2025.07.01-x86_64.tar.zst`.

## 🌟 Join the Adventure: Contribute!
We’re excited to build `arch-sandbox` with you! Whether you’re a Go developer, Linux guru, or UI enthusiast, your ideas are welcome. Here’s how you can contribute:
- **Snapshots**: Implement `snapshot.go` to save/restore sandbox states.
- **Graphs & Visuals**: Create dynamic graphs (e.g., Chart.js) for setup times or disk usage.
- **New Features**: Add commands to list, pause, or monitor sandboxes.
- **UI Enhancements**: Build a terminal progress bar or web dashboard.

**Get Started**:
1. Fork the repository.
2. Create a branch: `git checkout -b feature/your-awesome-idea`.
3. Commit: `git commit -m 'Add awesome idea'`.
4. Push: `git push origin feature/your-awesome-idea`.
5. Open a pull request.

Share ideas or report bugs via [issues](https://github.com/OminduD/arch-sandbox/issues) or join our [community](#) *(Placeholder: Add Discord or forum link)*. Let’s make `arch-sandbox` epic together! 🚀

## 📜 License
MIT License. See [LICENSE](LICENSE) for details.

---
