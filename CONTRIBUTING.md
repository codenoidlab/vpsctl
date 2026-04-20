# Contributing to VPSctl

Thank you for your interest in contributing to VPSctl! We welcome contributions from everyone, whether it's bug reports, feature requests, or code improvements.

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.22+** — Download from [golang.org](https://golang.org)
- **Git** — For version control
- **Make** — For build commands (pre-installed on macOS/Linux; use `choco install make` on Windows)
- **Basic Go knowledge** — Familiarity with Go syntax and conventions

### Setting Up Your Development Environment

```bash
# 1. Fork the repository on GitHub
# Visit: https://github.com/codenoidlab/vpsctl

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/vpsctl.git
cd vpsctl

# 3. Add upstream remote
git remote add upstream https://github.com/codenoidlab/vpsctl.git

# 4. Verify setup
git remote -v
# Should show:
#   origin    https://github.com/YOUR_USERNAME/vpsctl.git
#   upstream  https://github.com/codenoidlab/vpsctl.git

# 5. Install dependencies
go mod download

# 6. Build the project
make build

# 7. Run tests
go test ./...

# 8. Verify linting works
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

---

## 📋 Types of Contributions

### 🐛 Bug Reports

Found a bug? Please help us by reporting it!

**Before submitting:**
1. Check [existing issues](https://github.com/codenoidlab/vpsctl/issues)
2. Ensure the bug is reproducible
3. Gather debug information

**Submit here:** [Submit a bug report](https://github.com/codenoidlab/vpsctl/issues/new?template=bug_report.md)

**Include:**
```markdown
## Description
Brief description of the bug

## Reproduction Steps
1. ...
2. ...
3. ...

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: Ubuntu 22.04
- Go version: 1.22
- VPSctl version: 0.1.0

## Logs/Output
```
error output here
```

## Screenshots
(if applicable)
```

### 💡 Feature Requests

Have an idea? We'd love to hear it!

**Submit here:** [Submit a feature request](https://github.com/codenoidlab/vpsctl/issues/new?template=feature_request.md)

**Include:**
```markdown
## Description
What should VPSctl do?

## Use Case
Why would this be useful?

## Example
Show a mockup or use case

## Alternatives Considered
Other approaches?
```

### 🔧 Code Contributions

Want to write code? Here's how:

1. **Pick an issue** from [GitHub Issues](https://github.com/codenoidlab/vpsctl/issues)
   - Look for: `good first issue`, `help wanted`, `enhancement`
   - Comment: "I'd like to work on this"

2. **Follow Git workflow** (see [GIT_WORKFLOW.md](GIT_WORKFLOW.md))

3. **Create feature branch:**
   ```bash
   git checkout develop
   git pull upstream develop
   git checkout -b feature/123-your-feature-name
   ```

4. **Make your changes:**
   ```bash
   # Edit files
   git add .
   git commit -m "feat: Add awesome feature"
   ```

5. **Keep in sync:**
   ```bash
   git fetch upstream
   git rebase upstream/develop
   ```

6. **Push to your fork:**
   ```bash
   git push origin feature/123-your-feature-name
   ```

7. **Create a Pull Request** on GitHub

---

## 💻 Development Guidelines

### Code Style

VPSctl follows standard Go conventions:

```bash
# Format code (always do this before committing)
go fmt ./...

# Or use goimports (auto-organize imports)
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .

# Lint with golangci-lint
golangci-lint run
```

**Key principles:**
- Keep functions small and focused
- Use clear, descriptive variable names
- Add comments for exported functions
- No global state — use dependency injection
- All shell commands through `core.RunCommand()`

### Commit Messages

Follow the Conventional Commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example:**
```
feat(dashboard): Add CPU temperature graph

Implement a new graph widget showing CPU temperature
trends over the last 24 hours. Uses the new thermal
data provider from osadapter.

Fixes #456
Relates to #789
```

**Types:**
- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `style:` — Code style (formatting)
- `refactor:` — Code refactoring
- `perf:` — Performance improvement
- `test:` — Tests
- `chore:` — Build, dependencies, etc.

**Scopes (examples):**
- `dashboard` — Dashboard module
- `files` — File manager
- `node` — Node.js/PM2 management
- `git` — Git management
- `security` — Security module
- `ui` — UI/TUI framework
- `core` — Core functionality

### Testing

Always add tests for new features:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./internal/modules/dashboard

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
```

**Test file naming:**
- `module_test.go` — Tests for `module.go`
- Keep tests in same package (`package_test` vs `package`)

**Example test:**
```go
package dashboard_test

import (
	"testing"
	"github.com/vpsmanager/vps/internal/modules/dashboard"
)

func TestDataLoading(t *testing.T) {
	// Arrange
	model := dashboard.New()
	
	// Act
	// result := model.LoadData()
	
	// Assert
	// if result.Err != nil {
	//     t.Fatalf("Expected no error, got %v", result.Err)
	// }
}
```

### Adding a New Module

See [README.md](README.md#-development) for full instructions, but quick version:

1. Create `internal/modules/yourmodule/yourmodule.go`
2. Add Screen constant in `internal/tui/nav.go`
3. Add NavItem entry in same file
4. Wire it into `internal/tui/app.go` (4 places)
5. Implement `Model`, `New()`, `LoadCmd()`, `Update()`, `View()`
6. Test locally
7. Create PR

### Documentation

- Update [README.md](README.md) if adding new features
- Add comments to exported functions/types
- Update [GIT_WORKFLOW.md](GIT_WORKFLOW.md) if changing workflow
- Add examples in code comments
- Create docs for complex logic

---

## 📝 Pull Request Process

### Before Submitting

- [ ] Code builds: `make build`
- [ ] Tests pass: `go test ./...`
- [ ] Linter passes: `golangci-lint run`
- [ ] Fmt passes: `go fmt ./...`
- [ ] Security scan passes: `gosec ./...`
- [ ] No hardcoded values or secrets
- [ ] Commits are squashed and clean
- [ ] Branch is up to date with `develop`

### Creating the PR

1. **Go to GitHub:** [github.com/codenoidlab/vpsctl/pulls](https://github.com/codenoidlab/vpsctl/pulls)
2. **Click "New Pull Request"**
3. **Set base:** `develop` (not `main`)
4. **Fill in the template:**

```markdown
## Description
Brief description of changes

## Related Issue
Fixes #456

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Changes Made
- List specific changes
- Be detailed
- Include what was added/modified/removed

## Testing
- [ ] Tested locally
- [ ] Tested over SSH
- [ ] No breaking changes
- [ ] Backward compatible

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All checks passing
- [ ] No new warnings

## Screenshots (if applicable)
Attach screenshots for UI changes
```

### During Review

- Be responsive to feedback
- Discuss disagreements respectfully
- Make requested changes promptly
- Push new commits (don't force push after review starts)
- Request re-review when changes are made

### After Approval

Maintainers will merge your PR. We typically squash commits for clean history.

---

## 🚀 Development Workflow Examples

### Fix a Bug

```bash
# 1. Create bugfix branch
git checkout develop
git pull upstream develop
git checkout -b bugfix/123-fix-cpu-calc

# 2. Fix the bug
# ... edit files ...
go test ./...  # Verify fix

# 3. Commit
git add .
git commit -m "fix(dashboard): Correct CPU percentage calculation"

# 4. Push
git push origin bugfix/123-fix-cpu-calc

# 5. Create PR to develop
# ... on GitHub, create PR ...
```

### Add a Feature

```bash
# 1. Create feature branch
git checkout develop
git pull upstream develop
git checkout -b feature/456-docker-support

# 2. Implement feature
# ... create module ...
# ... write tests ...
go test ./...
golangci-lint run

# 3. Commit (multiple commits are ok)
git add internal/modules/docker/
git commit -m "feat(docker): Add Docker module skeleton"
git add internal/modules/docker/docker.go
git commit -m "feat(docker): Implement Docker container listing"

# 4. Before pushing, squash commits
git rebase -i upstream/develop
# Mark first commit as 'pick', rest as 'squash'

# 5. Push
git push origin feature/456-docker-support --force-with-lease

# 6. Create PR to develop
```

### Sync with Upstream

```bash
# 1. Fetch latest
git fetch upstream

# 2. Rebase on develop
git rebase upstream/develop

# 3. Force push (only if no one else is working on branch)
git push origin feature/branch --force-with-lease
```

---

## ❓ Getting Help

- **Questions:** [GitHub Discussions](https://github.com/codenoidlab/vpsctl/discussions)
- **Issues:** [GitHub Issues](https://github.com/codenoidlab/vpsctl/issues)
- **Email:** hello@codenoid.com
- **Chat:** Coming soon

---

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Our Git Workflow](GIT_WORKFLOW.md)

---

## 🙏 Code of Conduct

Please be respectful and inclusive. We follow the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

---

## 🎉 Thank You!

Every contribution helps make VPSctl better. Whether it's code, bug reports, documentation, or just spreading the word — **you're awesome!** 🚀

---

**Happy contributing!**
