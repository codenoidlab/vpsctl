# VPSctl Git Branching Strategy & Workflow

> Professional Git flow for a global-level open-source package repository

---

## 📌 Overview

This document defines the Git branching strategy, naming conventions, and workflow policies for **VPSctl** — a production-ready, globally-distributed open-source tool.

**Repository:** `github.com/codenoidlab/vpsctl`

---

## 🌳 Branch Model: Git Flow + Trunk-Based Development Hybrid

We use a **simplified Git Flow** optimized for open-source projects:

```
main (production/stable)
  ↑
  └─ release/* (release preparation)
     ↑
develop (integration branch)
  ↑
  ├─ feature/* (feature development)
  ├─ bugfix/* (bug fixes)
  ├─ hotfix/* (emergency production fixes)
  └─ docs/* (documentation updates)
```

### Branch Hierarchy

| Branch Type | Born From | Merges To | Lifecycle | Use Case |
|-------------|-----------|-----------|-----------|----------|
| `main` | N/A | N/A | ♾️ Permanent | **Production-ready code** |
| `develop` | `main` | `main` | ♾️ Permanent | **Integration & pre-release** |
| `feature/*` | `develop` | `develop` | 🔄 Temporary | New features / features in progress |
| `bugfix/*` | `develop` | `develop` | 🔄 Temporary | Bug fixes (non-critical) |
| `hotfix/*` | `main` | `main` + `develop` | 🔄 Temporary | Critical production fixes |
| `release/*` | `develop` | `main` | 🔄 Temporary | Release candidates & final prep |
| `docs/*` | `develop` | `develop` | 🔄 Temporary | Documentation-only changes |

---

## 📝 Branch Naming Conventions

### Format

```
<type>/<issue-number>-<short-description>
```

### Examples

**Feature branches:**
```
feature/123-dashboard-refresh
feature/456-nginx-ssl-config
feature/789-git-repo-search
```

**Bugfix branches:**
```
bugfix/42-cpu-calculation-fix
bugfix/101-file-permission-error
bugfix/202-ssh-timeout-issue
```

**Hotfix branches:**
```
hotfix/critical-security-patch
hotfix/deployment-failure-fix
```

**Release branches:**
```
release/0.2.0
release/1.0.0-rc1
```

**Documentation branches:**
```
docs/api-documentation
docs/installation-guide
docs/contributing-guide
```

### Naming Rules

✅ **DO:**
- Use lowercase letters
- Use hyphens to separate words
- Include GitHub issue number (if applicable)
- Keep names short but descriptive (max 50 chars)
- Use meaningful verbs: `add-`, `fix-`, `improve-`, `refactor-`, `update-`

❌ **DON'T:**
- Use spaces or underscores (use hyphens)
- Use UPPERCASE
- Use vague names like `bugfix/stuff` or `feature/work`
- Include version numbers in features (only in release branches)
- Use multiple slashes (`feature/sub/path` - one slash max)

### Issue Number Mapping

If using GitHub issues:
```
feature/123-feature-name  ← references GitHub issue #123
```

If no issue exists:
```
feature/anonymous-short-description  ← for small/urgent changes
```

---

## 🔄 Git Workflow: Step-by-Step

### 1️⃣ Starting a New Feature

**Scenario:** You're adding a new feature (e.g., Docker management screen)

```bash
# Ensure you're up to date with develop
git checkout develop
git pull origin develop

# Create feature branch from develop
git checkout -b feature/789-docker-management

# Start coding...
git add .
git commit -m "feat: Add Docker container management screen"
```

**Commit Message Format:**
```
<type>: <short description>

<optional detailed description>

Fixes #789
```

**Types:**
- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `style:` — Code style (formatting, missing semicolons, etc.)
- `refactor:` — Code refactoring
- `perf:` — Performance improvements
- `test:` — Adding tests
- `chore:` — Build process, dependencies, etc.

**Example commits:**
```
feat: Add Docker container management screen
fix: Correct CPU usage calculation in Dashboard
docs: Update installation guide for Go 1.22+
refactor: Extract command execution to separate module
```

### 2️⃣ Pushing to Remote

```bash
# Push to your fork (for external contributors)
git push origin feature/789-docker-management

# Or push directly to main repo (for maintainers)
git push origin feature/789-docker-management
```

### 3️⃣ Creating a Pull Request

**On GitHub:**
1. Go to `github.com/codenoidlab/vpsctl`
2. Click "New Pull Request"
3. Set:
   - **Base branch:** `develop`
   - **Compare branch:** `feature/789-docker-management`

**PR Title:**
```
[FEATURE] Docker Management Screen

or

Add Docker container management and monitoring
```

**PR Description Template:**
```markdown
## Description
Brief description of changes

## Related Issue
Fixes #789

## Changes Made
- Added Docker screen module
- Implemented container list and status
- Added container start/stop/restart actions

## Testing
- [ ] Tested locally
- [ ] Tested over SSH
- [ ] No breaking changes

## Screenshots (if applicable)
[attach screenshots]

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] No new warnings generated
```

### 4️⃣ Code Review & Merging

**Reviewer Checklist:**
- ✅ Code quality & style
- ✅ No hardcoded values
- ✅ Proper error handling
- ✅ Tests passing
- ✅ Documentation updated
- ✅ No security issues

**Merge Options:**

```bash
# Option A: Squash commits (clean history)
git checkout develop
git pull origin develop
git merge --squash feature/789-docker-management
git commit -m "feat: Add Docker container management screen"
git push origin develop

# Option B: Create merge commit (preserves history)
git checkout develop
git pull origin develop
git merge --no-ff feature/789-docker-management
git push origin develop

# Option C: Rebase (linear history)
git checkout feature/789-docker-management
git rebase develop
git checkout develop
git merge --ff-only feature/789-docker-management
git push origin develop
```

**Policy:** Use **squash + merge** for clarity, or **merge commit** for complex features.

### 5️⃣ Cleanup

```bash
# Delete local branch
git branch -d feature/789-docker-management

# Delete remote branch
git push origin --delete feature/789-docker-management

# Or delete on GitHub after PR is merged (auto-cleanup option)
```

---

## 🐛 Bug Fixing Workflow

### Non-Critical Bug (from `develop`)

```bash
# Create bugfix branch
git checkout develop
git pull origin develop
git checkout -b bugfix/42-cpu-calculation

# Fix and commit
git add .
git commit -m "fix: Correct CPU percentage calculation"
git push origin bugfix/42-cpu-calculation

# Create PR to develop
# Merge to develop
```

### Critical Production Bug (Hotfix from `main`)

```bash
# Create hotfix branch directly from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-security-patch

# Fix and commit
git add .
git commit -m "fix: Security patch for command injection vulnerability"
git push origin hotfix/critical-security-patch

# Create PR to main (urgent review needed)
# After merge to main, also merge back to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop
```

---

## 📦 Release Workflow

### Creating a Release (v0.2.0)

```bash
# Ensure develop is ready
git checkout develop
git pull origin develop

# Create release branch
git checkout -b release/0.2.0

# Update version numbers
# File: go.mod (if needed)
# File: internal/tui/app.go (if version string exists)
# Example:
#   const Version = "0.2.0"

git add .
git commit -m "chore: Bump version to 0.2.0"

# Create release PR to main for final review
git push origin release/0.2.0
```

**Release Checklist:**
```markdown
## Release v0.2.0 Checklist

- [ ] All features merged to `develop`
- [ ] All tests passing
- [ ] Version bumped in code
- [ ] CHANGELOG.md updated
- [ ] README.md reviewed
- [ ] Go dependencies up to date (`go mod tidy`)
- [ ] Build successful (`make build`, `make linux`)
- [ ] Documentation updated
- [ ] No security vulnerabilities
```

### Finishing the Release

```bash
# After release PR approved:
git checkout main
git pull origin main
git merge release/0.2.0
git tag -a v0.2.0 -m "Release version 0.2.0"
git push origin main
git push origin v0.2.0

# Merge back to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop

# Delete release branch
git branch -d release/0.2.0
git push origin --delete release/0.2.0
```

---

## 👥 Contributor Workflow (External)

### For Open-Source Contributors

```bash
# 1. Fork repository on GitHub
# (creates yourname/vpsctl)

# 2. Clone your fork
git clone https://github.com/yourname/vpsctl.git
cd vpsctl

# 3. Add upstream remote
git remote add upstream https://github.com/codenoidlab/vpsctl.git

# 4. Create feature branch from develop
git fetch upstream
git checkout -b feature/123-my-feature upstream/develop

# 5. Make changes
git add .
git commit -m "feat: Add my feature"

# 6. Keep in sync with upstream
git fetch upstream
git rebase upstream/develop

# 7. Push to your fork
git push origin feature/123-my-feature

# 8. Create Pull Request on GitHub
# (GitHub will detect your fork and suggest PR)

# 9. After merge, cleanup
git checkout develop
git pull upstream develop
git branch -d feature/123-my-feature
```

---

## 🔐 Protected Branches Configuration

### GitHub Settings (Repo → Settings → Branches)

**Protect `main` branch:**
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass (CI/CD)
- ✅ Require branches to be up to date
- ✅ Include administrators
- ✅ Restrict who can push to matching branches

**Protect `develop` branch:**
- ✅ Require pull request reviews before merging (1 reviewer minimum)
- ✅ Require status checks to pass (CI/CD)
- ⚪ Allow force pushes (optional, not recommended)
- ⚪ Allow deletions (should be disabled)

### Required Status Checks

Configure these in GitHub:
```
- build (Go compilation)
- test (go test ./...)
- lint (golangci-lint or similar)
- coverage (code coverage threshold)
```

---

## 📊 Commit History Rules

### Commit Message Policy

**Required Format:**
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example:**
```
feat(dashboard): Add real-time CPU graph

Implement a new graph widget for CPU usage trends
showing 24-hour historical data. Uses charmbracelet's
plot library for terminal rendering.

Fixes #456
Relates to #789
```

### Squash Commits Before Merging

❌ **Bad:**
```
Commit 1: WIP - dashboard update
Commit 2: Fix typo in dashboard
Commit 3: Revert previous changes
Commit 4: Final dashboard feature
```

✅ **Good:**
```
Commit 1: feat: Add real-time CPU graph to dashboard
```

**How to squash:**
```bash
git rebase -i HEAD~4  # Interactive rebase for 4 commits
# Mark first commit as 'pick', rest as 'squash'
# Edit combined message
git push origin feature/branch --force-with-lease
```

---

## 🚦 CI/CD Pipeline Requirements

Each push triggers:

1. **Build:** `go build`
2. **Test:** `go test ./...`
3. **Lint:** `golangci-lint run`
4. **Coverage:** Report test coverage > 70%
5. **Security:** `gosec ./...`
6. **Documentation:** Check for doc updates

**PR cannot be merged if any check fails.**

---

## 📋 Common Scenarios

### Scenario 1: Adding a Feature

```bash
git checkout develop && git pull origin develop
git checkout -b feature/123-new-screen
# ... make changes ...
git add .
git commit -m "feat: Add new management screen"
git push origin feature/123-new-screen
# Create PR to develop
# After review and approval, merge via GitHub
git fetch origin && git branch -d feature/123-new-screen
```

### Scenario 2: Fixing a Bug

```bash
git checkout develop && git pull origin develop
git checkout -b bugfix/42-fix-issue
# ... fix bug ...
git add .
git commit -m "fix: Resolve issue in file manager"
git push origin bugfix/42-fix-issue
# Create PR to develop
# After merge, cleanup
```

### Scenario 3: Emergency Production Fix

```bash
git checkout main && git pull origin main
git checkout -b hotfix/critical-fix
# ... fix critical issue ...
git add .
git commit -m "fix: Critical security patch"
git push origin hotfix/critical-fix
# Create PR to main with urgent flag
# After merge to main, also merge to develop
git checkout develop && git pull origin develop
git merge main && git push origin develop
# Cleanup
git branch -d hotfix/critical-fix
```

### Scenario 4: Releasing a Version

```bash
git checkout develop && git pull origin develop
git checkout -b release/v0.3.0
# Update version numbers, CHANGELOG
git add .
git commit -m "chore: Prepare v0.3.0 release"
git push origin release/v0.3.0
# Create PR to main for final review
# After approval, merge and tag
git checkout main && git merge release/v0.3.0
git tag -a v0.3.0 -m "Release v0.3.0"
git push origin main v0.3.0
# Merge back to develop
git checkout develop && git merge main && git push origin develop
# Cleanup
git branch -d release/v0.3.0
```

---

## 🔄 Pull & Push Configuration

### Git Configuration for Developers

```bash
# Global config (first-time setup)
git config --global user.name "Your Name"
git config --global user.email "your@email.com"

# Project-specific config (optional)
cd vpsctl
git config user.name "Your Name"
git config user.email "your@email.com"

# Set default branch strategy
git config pull.rebase false  # Use merge (default)
# or
git config pull.rebase true   # Use rebase

# Set default push behavior
git config push.default current  # Push current branch with same name
git config push.followTags true  # Auto-push tags
```

### .gitconfig for Convenience

Add to `~/.gitconfig`:

```ini
[user]
    name = Your Name
    email = your@email.com

[core]
    editor = vim  # or nano, code, etc.
    autocrlf = false

[pull]
    rebase = false

[push]
    default = current
    followTags = true

[alias]
    # Common aliases
    st = status
    co = checkout
    br = branch
    ci = commit
    amend = commit --amend --no-edit
    undo = reset --soft HEAD~1
    logs = log --oneline --graph --all
    sync = !git fetch origin && git rebase origin/develop

[merge]
    conflictstyle = zdiff3

[rebase]
    autosquash = true
```

### Push & Pull Best Practices

```bash
# Before pushing, always pull latest
git pull origin feature/branch

# Push with lease (safer than force)
git push origin feature/branch --force-with-lease

# Pull with specific strategy
git pull origin develop --rebase  # Rebase pull
git pull origin develop --no-rebase  # Merge pull

# Keep fork in sync with upstream
git fetch upstream
git rebase upstream/develop
```

---

## 🚀 Requirements Checklist

### Before Committing

- [ ] Code follows Go conventions
- [ ] No debug statements or console.log equivalents
- [ ] All tests pass locally: `go test ./...`
- [ ] Linter passes: `golangci-lint run`
- [ ] Security check passes: `gosec ./...`
- [ ] Code builds: `make build`

### Before Pushing

- [ ] Commit message is clear and follows format
- [ ] Branch is up to date: `git pull --rebase origin develop`
- [ ] No sensitive data (API keys, passwords) in code
- [ ] Documentation is updated (README, comments, docs/)

### Before Creating PR

- [ ] PR title is descriptive
- [ ] PR description explains changes
- [ ] Related issue is linked
- [ ] Screenshots/demos included (if UI change)
- [ ] Commit history is clean (squashed if needed)
- [ ] All checks pass (CI/CD green)

### Before Merging PR

- [ ] At least 1 code review approval
- [ ] All discussions resolved
- [ ] All checks passing
- [ ] Commits are squashed
- [ ] Branch is up to date with target

---

## 📚 Useful Git Commands

```bash
# View status clearly
git status

# See commit history
git log --oneline --graph --all

# See who changed a line
git blame filename.go

# Find a commit
git log --grep="search term" --oneline

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1

# Revert a commit (safe for public branches)
git revert <commit-hash>

# Cherry-pick a commit to another branch
git cherry-pick <commit-hash>

# Create a patch file
git format-patch origin/develop -o patches/

# Apply a patch
git apply patches/0001-fix-something.patch

# Stash changes temporarily
git stash
git stash pop

# Check what's ahead/behind
git rev-list --count origin/develop..develop  # ahead
git rev-list --count develop..origin/develop  # behind
```

---

## 🔍 Example Repository Structure Over Time

```
main
  ├─ v0.1.0  ← Initial release
  ├─ v0.2.0  ← Feature release
  └─ v0.3.0  ← Feature release

develop
  ├─ feature/456-docker-support (being reviewed)
  ├─ feature/789-database-screen (in progress)
  └─ docs/api-reference (ready to merge)

feature/*
  ├─ feature/456-docker-support
  ├─ feature/500-kubernetes-integration
  └─ feature/789-database-screen

hotfix/*
  └─ hotfix/security-patch

release/*
  └─ release/0.3.0
```

---

## ✅ Team Standards Summary

| Aspect | Standard |
|--------|----------|
| **Main Branch** | `main` — Production only |
| **Development** | `develop` — Next release |
| **Feature Branch** | `feature/issue-name` from `develop` |
| **Bug Branch** | `bugfix/issue-name` from `develop` |
| **Hotfix Branch** | `hotfix/issue-name` from `main` |
| **Release Branch** | `release/version` from `develop` |
| **Commit Message** | `<type>: <description>` |
| **Merge Strategy** | Squash + merge (clean history) |
| **PR Reviews** | 1+ approval required |
| **CI/CD** | All checks must pass |
| **Push/Pull** | `--force-with-lease` only |
| **Tag Format** | `v0.1.0` (semantic versioning) |

---

## 🎯 Getting Started

### Project Maintainers

```bash
# Initial setup
git clone https://github.com/codenoidlab/vpsctl.git
cd vpsctl
git checkout develop

# Create feature / hotfix as needed
git checkout -b feature/123-feature-name
# ... work ...
git push origin feature/123-feature-name
```

### External Contributors

```bash
# Fork and setup
git clone https://github.com/yourname/vpsctl.git
cd vpsctl
git remote add upstream https://github.com/codenoidlab/vpsctl.git

# Create feature from upstream develop
git checkout -b feature/123-feature upstream/develop
# ... work ...
git push origin feature/123-feature
# Create PR on GitHub
```

---

## 📖 Additional Resources

- [Git Documentation](https://git-scm.com/doc)
- [GitHub Flow Guide](https://guides.github.com/introduction/flow/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [CocoaPods Git Flow](https://guides.cocoapods.org/using/pods-and-git.html)

---

## 📞 Questions?

- **Issues:** [github.com/codenoidlab/vpsctl/issues](https://github.com/codenoidlab/vpsctl/issues)
- **Discussions:** [github.com/codenoidlab/vpsctl/discussions](https://github.com/codenoidlab/vpsctl/discussions)
- **Email:** hello@codenoid.com

---

**Version:** 1.0  
**Last Updated:** April 2026  
**Maintained by:** Codenoid Lab
