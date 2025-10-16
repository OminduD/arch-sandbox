<div align="center">

```
    ___             __       _____                ____              
   /   |  __________/ /_    / ___/____ _____  ___/ / /_  ____  _  __
  / /| | / ___/ ___/ __ \   \__ \/ __ `/ __ \/ __  / __ \/ __ \| |/_/
 / ___ |/ /  / /__/ / / /  ___/ / /_/ / / / / /_/ / /_/ / /_/ />  <  
/_/  |_/_/   \___/_/ /_/  /____/\__,_/_/ /_/\__,_/_.___/\____/_/|_|  
```

# 🏝️ Arch-Sandbox

### Create and manage isolated Arch Linux environments with ease

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?logo=go)](https://go.dev/)
[![Arch Linux](https://img.shields.io/badge/Arch-Linux-1793D1?logo=arch-linux)](https://archlinux.org/)

**arch-sandbox** is a powerful command-line tool that spins up isolated Arch Linux environments using overlay filesystems and `systemd-nspawn`. Perfect for developers, system administrators, and Linux enthusiasts who need safe, disposable environments for testing and experimentation.

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Contributing](#-contributing) • [License](#-license)

</div>

---

## ✨ Features

- 🔒 **Isolated Environments** - Run Arch Linux sandboxes without affecting your host system
- 📦 **Overlay Filesystem** - Keep changes separate using overlayfs technology
- 💾 **Persistent or Ephemeral** - Choose to keep or discard your sandbox after use
- 📸 **Snapshot Management** - Save and restore sandbox states
- 📦 **Package Installation** - Install packages directly in persistent sandboxes
- 🎯 **Simple CLI** - Powered by Cobra for an intuitive command-line experience
- ⚡ **Lightweight** - Minimal overhead, maximum performance

## 📋 Prerequisites

Ensure you have the following installed on your Arch Linux system:
- **systemd-nspawn** - for containerization
- **mount** - for overlay filesystem operations
- **pacman** - for Arch Linux package management
- **zstd** - for `.tar.zst` archive handling

Install prerequisites on Arch Linux:
```bash
sudo pacman -S systemd zstd
```

> **Note:** You'll need root privileges (`sudo`) to run arch-sandbox due to `systemd-nspawn` and `mount` requirements.

## 🚀 Installation

### Method 1: Install from AUR (Recommended)

The easiest way to install arch-sandbox is through the Arch User Repository (AUR):

```bash
# Using yay
yay -S arch-sandbox

# Or using paru
paru -S arch-sandbox
```

> **Note:** The AUR package may be in the process of being published. If it's not available yet, please use Method 2 (Manual Installation) below.

This method automatically handles dependencies and keeps the tool up-to-date with your system.

### Method 2: Manual Installation from Source

If you prefer to build from source or want the latest development version:

#### Step 1: Install Build Dependencies
```bash
sudo pacman -S go git systemd zstd
```

#### Step 2: Clone the Repository
```bash
git clone https://github.com/OminduD/arch-sandbox.git
cd arch-sandbox
```

#### Step 3: Build the Binary
```bash
go build -o arch-sandbox
```

#### Step 4: Install Globally (Optional)
```bash
sudo install -Dm755 arch-sandbox /usr/local/bin/arch-sandbox
```

Or copy it to a directory in your PATH:
```bash
sudo mv arch-sandbox /usr/local/bin/
```

#### Verify Installation
```bash
arch-sandbox --help
```

## 🛠️ Usage

### Quick Start

Create a disposable sandbox for quick testing:
```bash
sudo arch-sandbox new mysandbox
```

Create a persistent sandbox that survives reboots:
```bash
sudo arch-sandbox new mysandbox --persist
```

### Available Commands

#### Create a New Sandbox
```bash
arch-sandbox new <name> [flags]
```

**Flags:**
- `-p, --persist` - Keep the sandbox after exiting (default: cleanup on exit)
- `--base-dir string` - Base directory for sandboxes (default: `~/.arch-sandbox`)

**Examples:**
```bash
# Create a disposable sandbox
sudo arch-sandbox new testbox

# Create a persistent sandbox
sudo arch-sandbox new devbox --persist

# Create a sandbox in a custom location
sudo arch-sandbox new projectbox --persist --base-dir /data/sandboxes
```

#### Install Packages
Install packages directly in a persistent sandbox:
```bash
sudo arch-sandbox install <sandbox-name> <package>
```

**Example:**
```bash
# Install vim in a sandbox
sudo arch-sandbox install devbox vim

# Install git in a sandbox
sudo arch-sandbox install devbox git
```

#### Manage Snapshots
Save and restore sandbox states:
```bash
# Save a snapshot
sudo arch-sandbox snapshot <sandbox-name> save <snapshot-id>

# Restore a snapshot
sudo arch-sandbox snapshot <sandbox-name> restore <snapshot-id>

# List snapshots
sudo arch-sandbox snapshot <sandbox-name> list
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

### Inside the Sandbox
You'll enter a `/bin/bash` shell where you can:
- 📦 Install packages with `pacman`
- 🧪 Test scripts or configurations safely
- 🔧 Experiment with system changes
- 📝 Develop and test software
- 🚪 Exit with `exit` or `Ctrl+D`

### Managing Sandboxes

#### Cleanup
- **Disposable Sandboxes**: Automatically deleted on exit
- **Persistent Sandboxes**: Stored in `~/.arch-sandbox/<name>`

Delete a persistent sandbox manually:
```bash
rm -rf ~/.arch-sandbox/<name>
```

#### List All Sandboxes
```bash
ls -la ~/.arch-sandbox/
```

## ⚠️ Important Notes
- 🔐 Run as `root` or with `sudo` for `systemd-nspawn` and `mount` operations
- 🌐 Internet access is required for tarball download
- 📦 Tarball source: `https://archive.archlinux.org/iso/2025.07.01/archlinux-bootstrap-2025.07.01-x86_64.tar.zst`

## 🤝 Contributing

We welcome contributions from everyone! Whether you're fixing bugs, adding features, or improving documentation, your help makes arch-sandbox better.

### How to Contribute

1. **🍴 Fork the Repository**
   ```bash
   # Click the "Fork" button on GitHub, then clone your fork
   git clone https://github.com/YOUR-USERNAME/arch-sandbox.git
   cd arch-sandbox
   ```

2. **🌿 Create a Feature Branch**
   ```bash
   git checkout -b feature/amazing-feature
   # or
   git checkout -b fix/bug-description
   ```

3. **💻 Make Your Changes**
   - Write clean, readable code
   - Follow existing code style and conventions
   - Add comments for complex logic
   - Test your changes thoroughly

4. **✅ Test Your Changes**
   ```bash
   # Build the project
   go build -o arch-sandbox
   
   # Test the binary
   sudo ./arch-sandbox new testbox
   ```

5. **📝 Commit Your Changes**
   ```bash
   git add .
   git commit -m "Add: Brief description of your changes"
   ```
   
   **Commit Message Guidelines:**
   - `Add:` for new features
   - `Fix:` for bug fixes
   - `Update:` for updates to existing features
   - `Docs:` for documentation changes

6. **🚀 Push to Your Fork**
   ```bash
   git push origin feature/amazing-feature
   ```

7. **🎯 Open a Pull Request**
   - Go to the original repository on GitHub
   - Click "New Pull Request"
   - Select your fork and branch
   - Describe your changes in detail
   - Link any related issues

### Ideas for Contributions

Here are some areas where you can help:

- 🔧 **Features**
  - Add list command to show all sandboxes
  - Implement sandbox pause/resume functionality
  - Add resource monitoring (CPU, memory usage)
  - Create a web dashboard for managing sandboxes

- 📝 **Documentation**
  - Add more usage examples
  - Create video tutorials
  - Translate README to other languages
  - Write blog posts about use cases

- 🐛 **Bug Fixes**
  - Report bugs via [issues](https://github.com/OminduD/arch-sandbox/issues)
  - Fix existing issues
  - Improve error messages

- 🎨 **Design**
  - Create a logo or mascot
  - Design promotional graphics
  - Improve CLI output formatting

- 🧪 **Testing**
  - Add unit tests
  - Create integration tests
  - Test on different Arch Linux setups

### Development Setup

```bash
# Clone the repository
git clone https://github.com/OminduD/arch-sandbox.git
cd arch-sandbox

# Install dependencies
go mod download

# Build the project
go build -o arch-sandbox

# Run with sudo
sudo ./arch-sandbox --help
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions

### Getting Help

- 📫 Open an [issue](https://github.com/OminduD/arch-sandbox/issues) for questions
- 💬 Join discussions in the [Discussions](https://github.com/OminduD/arch-sandbox/discussions) tab
- 🐛 Report bugs with detailed reproduction steps

Let's make arch-sandbox amazing together! 🚀

## 📜 License
MIT License. See [LICENSE](LICENSE) for details.

---

<div align="center">
Made with ❤️ by the Arch-Sandbox community
</div>
