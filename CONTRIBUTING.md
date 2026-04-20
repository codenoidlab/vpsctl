# Contributing to VPSctl

First off, thank you for considering contributing to VPSctl! It is people like you that make open-source software such a great community. We welcome contributions from everyone, whether it's bug reports, feature requests, documentation improvements, or code changes.

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.22+** — Download from [golang.org](https://golang.org)
- **Git** — For version control
- **Make** — For build commands (pre-installed on macOS/Linux; use `choco install make` on Windows)
- **Basic Go knowledge** — Familiarity with Go syntax and BubbleTea TUI conventions

---

## 🤝 Standard Open Source Workflow

We follow the standard GitHub "Fork and Pull Request" workflow. Here is exactly how to get started:

### 1. Fork and Clone
1. Click the **Fork** button at the top-right of the VPSctl repository.
2. Clone your fork to your local machine:
   ```bash
   git clone https://github.com/YOUR-USERNAME/vpsctl.git
   cd vpsctl
   ```
3. Add the original repository as your `upstream` remote so you can stay up-to-date:
   ```bash
   git remote add upstream https://github.com/codenoidlab/vpsctl.git
   ```

### 2. Create a Branch
Always create a new branch for your work. Never push directly to your `main` or `develop` branch.
```bash
# First, ensure you have the latest code
git checkout develop
git pull upstream develop

# Then create your new branch (use feature/ or bugfix/ prefix)
git checkout -b feature/my-awesome-feature
```

### 3. Make Your Changes
Write your code, add tests, and update documentation. Be sure to test everything locally!

### 4. Commit Your Changes
We follow the **Conventional Commits** format for our commit messages:
```
<type>(<scope>): <subject>
```
**Examples:**
- `feat(dashboard): add CPU temperature graph`
- `fix(git): resolve nil pointer panic on empty repo`
- `docs: update setup instructions`

**Types to use:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.

```bash
git add .
git commit -m "feat(ui): add new responsive layout for tables"
```

### 5. Push and Open a Pull Request
Push your branch to your personal fork on GitHub:
```bash
git push origin feature/my-awesome-feature
```
Then, go to the original `codenoidlab/vpsctl` repository. You will see a prompt to **Compare & pull request**.
- Ensure you set the **base branch** to `develop`.
- Provide a clear title and description explaining what your PR solves or adds.

---

## 🛠️ Development Setup & Rules

To run VPSctl locally on your machine during development:

```bash
# 1. Install dependencies
go mod download

# 2. Format your code (Always do this before committing!)
go fmt ./...

# 3. Build and test
make build     # Builds the executable
go test ./...  # Runs tests

# 4. Run the application directly to test UI changes
go run ./cmd/vps
```

### Code Style Guidelines
- Keep functions small and focused.
- Use clear, descriptive variable names.
- No global state — use dependency injection.
- All shell commands must be executed safely through `core.RunCommand()`.
- Run `golangci-lint run` locally to ensure you don't introduce linter errors.

---

## 🐛 Found a Bug?
If you find a bug in the source code or a mistake in the documentation, please submit an issue to our [GitHub Repository](https://github.com/codenoidlab/vpsctl/issues). 

**Please include:**
- A clear description of the bug
- Exact steps to reproduce the issue
- Your OS, Go version, and VPSctl version
- Any relevant crash logs or screenshots

## 💡 Proposing a Feature
If you have an idea for a new feature, please submit an issue first to discuss it with the maintainers before you spend time coding it! This ensures your hard work aligns with the project's roadmap.

---

🎉 **Thank You!** Every contribution helps make VPSctl better! 🚀

