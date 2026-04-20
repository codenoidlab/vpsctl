# VPSctl

> **All-in-one terminal control panel for managing Ubuntu VPS servers**

[![GitHub Release](https://img.shields.io/github/v/release/codenoidlab/vpsctl?style=flat-square)](https://github.com/codenoidlab/vpsctl/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/codenoidlab/vpsctl/build.yml?style=flat-square)](https://github.com/codenoidlab/vpsctl/actions)

VPSctl is a **keyboard-driven terminal UI** for complete VPS management. SSH into your Ubuntu server, type `vpsctl`, and access a unified dashboard with 8 integrated screens for monitoring, administration, and service management.

**Website:** [vpsctl.codenoid.com](https://vpsctl.codenoid.com)  
**Repository:** [github.com/codenoidlab/vpsctl](https://github.com/codenoidlab/vpsctl)

---

## ✨ Features

### 🎛️ 8 Integrated Screens

| # | Screen | Capabilities |
|---|--------|--------------|
| **1** | **Dashboard** | CPU, RAM, disk monitoring + service status + uptime |
| **2** | **Files** | Browse, copy, move, delete files with full navigation |
| **3** | **Node / PM2** | Start, stop, restart Node apps + live logs |
| **4** | **Nginx** | Web server config + Cloudflare DNS integration |
| **5** | **Git** | Multi-repo management (pull, push, status) |
| **6** | **Packages** | apt, npm, snap package management |
| **7** | **Monitor** | Real-time process list + kill processes |
| **8** | **Security** | UFW firewall, SSH keys, Fail2ban config |

### ⚡ Quick Benefits

- **No Dependencies** — Single binary, works on any Ubuntu system
- **Fast & Lightweight** — Written in Go, instant startup
- **Keyboard-Driven** — Navigate everything without mouse
- **SSH-Native** — Works perfectly over SSH with low latency
- **Modular Design** — Easy to extend with new screens
- **Clean Exit** — Uses terminal alternate screen, leaves no artifacts

---

## 🚀 Quick Start

### Installation

#### Option 1: Download Pre-Built Binary
```bash
# Coming soon — check releases page
curl -L https://vpsctl.codenoid.com/install.sh | bash
```

#### Option 2: Build from Source
Requires **Go 1.22+**

```bash
# Clone repository
git clone https://github.com/codenoidlab/vpsctl.git
cd vpsctl

# Build for your OS
make build

# Install to system
make install
```

#### Option 3: Build for Remote Linux Server
```bash
# From Mac/Windows, cross-compile for Linux
make linux

# Copy to server
scp vpsctl-linux-amd64 user@yourserver:/tmp/

# SSH and install
ssh user@yourserver
cat /tmp/vpsctl-linux-amd64 > /usr/local/bin/vpsctl
chmod +x /usr/local/bin/vpsctl
```

### First Run

```bash
# Simply type:
vpsctl

# Or with version info:
vpsctl --version
```

---

## ⌨️ Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `1`–`8` | Jump directly to screen N |
| `Tab` / `Shift+Tab` | Cycle through screens |
| `↑` / `↓` or `k` / `j` | Navigate lists |
| `Enter` | Select/Execute |
| `r` | Refresh current screen |
| `/` | Search (module-specific) |
| `q` or `Ctrl+C` | Quit |

### Module-Specific

Each screen has additional shortcuts for module-specific actions (edit, delete, execute). Press `?` in any screen for help.

---

## 📋 System Requirements

- **OS:** Ubuntu 18.04 LTS or newer (or any Linux distribution with bash)
- **Go:** 1.22+ (only for building from source)
- **Root/Sudo:** Some operations (firewall, package install) require elevated privileges
- **SSH:** Works over SSH, local terminal, or tmux/screen
- **Terminal:** Any standard terminal emulator (no special font required)

---

## 🏗️ Project Structure

```
vpsctl/
├── cmd/
│   └── vps/
│       └── main.go              ← Entry point
├── core/
│   └── runner.go                ← Shared command execution layer
├── internal/
│   ├── modules/                 ← Feature screens
│   │   ├── dashboard/
│   │   ├── files/
│   │   ├── node/
│   │   ├── nginx/
│   │   ├── git/
│   │   ├── packages/
│   │   ├── monitor/
│   │   └── security/
│   ├── style/                   ← Centralized styling
│   ├── tui/                     ← UI framework & routing
│   │   ├── app.go
│   │   ├── nav.go
│   │   └── styles.go
├── pkg/
│   └── osadapter/               ← OS-level file reading
├── scripts/
│   └── install.sh               ← Installation script
├── Makefile                     ← Build commands
├── go.mod / go.sum              ← Go dependencies
└── README.md                    ← This file
```

---

## 🔧 Development

### Build Commands

```bash
# Build for current OS
make build

# Cross-compile for Linux amd64 (from Mac/Windows)
make linux

# Build both
make all

# Clean build artifacts
make clean

# Install to /usr/local/bin
make install
```

### Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy
```

### Adding a New Screen

1. **Add Screen constant** in `internal/tui/nav.go`:
   ```go
   const ScreenDatabase Screen = iota  // 9
   ```

2. **Add NavItem** in same file:
   ```go
   {"🗄", "Database", "9", ScreenDatabase},
   ```

3. **Create module** `internal/modules/database/database.go` with:
   - `Model` struct
   - `New()` constructor
   - `LoadCmd()` for data fetching
   - `Update()` for input handling
   - `View()` for rendering

4. **Wire into app** in `internal/tui/app.go`:
   - Import the module
   - Add field to `AppModel`
   - Add case in `loadCmd()` switch
   - Add case in `View()` switch

5. **Rebuild**:
   ```bash
   make build
   ```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep modules focused and single-purpose
- All shell commands must go through `core.RunCommand()`
- No global state — use dependency injection

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

---

## 🚢 Deployment

### Single Server

```bash
# Cross-compile
make linux

# Deploy
scp vpsctl-linux-amd64 user@server:/usr/local/bin/vpsctl
ssh user@server "chmod +x /usr/local/bin/vpsctl && vpsctl"
```

### Multiple Servers

```bash
# Build once
make linux

# Deploy to multiple servers
for server in server1 server2 server3; do
  scp vpsctl-linux-amd64 user@$server:/usr/local/bin/vpsctl
  ssh user@$server "chmod +x /usr/local/bin/vpsctl"
done
```

### Systemd Service (Optional)

Create `/etc/systemd/user/vpsctl.service`:

```ini
[Unit]
Description=VPSctl Server Manager
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/vpsctl
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
```

Then:
```bash
systemctl --user enable vpsctl
systemctl --user start vpsctl
```

---

## 📖 Architecture

### Design Principles

- **Modular** — Each screen is independent; add features without touching others
- **Centralized** — All shell commands go through `core.RunCommand()` for testing/logging
- **Reactive** — BubbleTea game loop: input → update state → render
- **Clean** — No global state; use dependency injection and message passing

### BubbleTea Architecture

VPSctl uses [BubbleTea](https://github.com/charmbracelet/bubbletea) for the TUI framework:

```
User Input
    ↓
BubbleTea Update()  ← Process keypress
    ↓
Module's Update()   ← Handle business logic
    ↓
View() → Terminal   ← Render screen
    ↓
Loop back to input
```

### Data Flow

```
Module needs OS data
    ↓
Run core.RunCommand() asynchronously
    ↓
Receive result as BubbleTea message
    ↓
Update() processes message
    ↓
View() renders new state
```

---

## 🤝 Contributing

We welcome contributions! Whether it's bug reports, feature requests, or code:

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### Development Guidelines

- Write clear commit messages
- Add tests for new features
- Update documentation
- Follow Go conventions
- Test on actual Ubuntu systems if possible

---

## 🐛 Troubleshooting

### Binary won't run

```bash
# Check permissions
ls -la /usr/local/bin/vpsctl

# Make executable
chmod +x /usr/local/bin/vpsctl

# Run with full path
/usr/local/bin/vpsctl
```

### Slow over SSH

- Ensure SSH connection has low latency
- Try native SSH (not through jumping hosts)
- Check server CPU load (use Dashboard screen)

### Module data not updating

- Press `r` to manually refresh current screen
- Check server permissions (some operations need sudo)
- View logs: check UFW, systemd for permission errors

### Build fails on macOS

```bash
# Ensure Go 1.22+ is installed
go version

# Update Go
brew upgrade go

# Try building
make build
```

---

## 📝 License

MIT License — See [LICENSE](LICENSE) file for details.

You're free to use, modify, and distribute VPSctl in commercial and personal projects.

---

## 📚 Resources

- **Go Documentation:** [golang.org](https://golang.org)
- **BubbleTea Guide:** [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- **Lipgloss Styling:** [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
- **DigitalOcean VPS Guide:** [docs.digitalocean.com](https://docs.digitalocean.com)

---

## 🙋 Support

- **Issues & Bugs:** [GitHub Issues](https://github.com/codenoidlab/vpsctl/issues)
- **Discussions:** [GitHub Discussions](https://github.com/codenoidlab/vpsctl/discussions)
- **Website:** [vpsctl.codenoid.com](https://vpsctl.codenoid.com)
- **Email:** support@codenoid.com

---

## 🎯 Roadmap

### v0.2.0 (Next)
- [ ] Command history & search
- [ ] Configuration file support
- [ ] SSH key generation wizard
- [ ] Enhanced error messages

### v0.3.0
- [ ] Docker container management
- [ ] Multiple VPS management (same UI)
- [ ] Alert system (CPU/disk warnings)
- [ ] Settings/preferences screen

### v1.0.0
- [ ] Plugin system
- [ ] Database management
- [ ] Web dashboard companion
- [ ] Commercial support

---

## 🏆 Credits

Built with ❤️ by [Codenoid Lab](https://codenoid.com)

**Special thanks to:**
- [Charmbracelet](https://charm.sh) — BubbleTea & Lipgloss
- [Go Community](https://golang.org) — Amazing language and ecosystem

---

## 📞 Connect

- **Website:** [vpsctl.codenoid.com](https://vpsctl.codenoid.com)
- **GitHub:** [github.com/codenoidlab/vpsctl](https://github.com/codenoidlab/vpsctl)
- **Twitter:** [@codenoidlab](https://twitter.com/codenoidlab)
- **Email:** hello@codenoid.com

---

**Star ⭐ this repository if you find it useful!**

Made with Go 🚀
